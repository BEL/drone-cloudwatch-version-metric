package main

import (
    "log"
    "fmt"
    "flag"
    "os"
    "strconv"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/awserr"
    "github.com/aws/aws-sdk-go/aws/awsutil"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/cloudwatch"
)

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}

func addMetric(name, unit string, value float64, dimensions []*cloudwatch.Dimension, metricData []*cloudwatch.MetricDatum) (ret []*cloudwatch.MetricDatum, err error) {
    _metric := cloudwatch.MetricDatum{
        MetricName: aws.String(name),
        Unit:       aws.String(unit),
        Value:      aws.Float64(value),
        Dimensions: dimensions,
    }
    metricData = append(metricData, &_metric)
    return metricData, nil
}

func putMetric(metricdata []*cloudwatch.MetricDatum, namespace, region string) error {

    session := session.New(&aws.Config{Region: &region})
    svc := cloudwatch.New(session)
    
    metric_input := &cloudwatch.PutMetricDataInput{
        MetricData: metricdata,
        Namespace:  aws.String(namespace),
    }

    resp, err := svc.PutMetricData(metric_input)
    if err != nil {
        if awsErr, ok := err.(awserr.Error); ok {
            return fmt.Errorf("[%s] %s", awsErr.Code, awsErr.Message)
        } else if err != nil {
            return err
        }
    }
    log.Println(awsutil.StringValue(resp))
    return nil
}

func main() {
    //Variable setting  
    var metricData []*cloudwatch.MetricDatum
    var ns = flag.String("namespace", "TestSpace", "CloudWatch metric namespace")//TODO:Change namespace to something meaningful.
    var dim []*cloudwatch.Dimension
    var err error

    //Dimensions of the metrica - Service name and environment
    var serviceName = "Service name"
    var serviceNameValue = getEnv("PLUGIN_SERVICE_NAME","my-service")
    var serviceEnvironment = "Environment"
    var serviceEnvironmentValue = getEnv("PLUGIN_SERVICE_ENV","test")

    deployment , err := strconv.ParseFloat(getEnv("PLUGIN_BUILD_NUMBER","0"),64)
    if err != nil {
         log.Fatal("Invalid build number:", err)
    }
    
    name := cloudwatch.Dimension{
            Name:  aws.String(serviceName),
            Value: aws.String(serviceNameValue),
        }
    env := cloudwatch.Dimension{
            Name:  aws.String(serviceEnvironment),
            Value: aws.String(serviceEnvironmentValue),
    }
    dim = append(dim,&name,&env)

    //Metric data addition
    metricData, err = addMetric("Deployment", "Count", deployment, dim, metricData)
        if err != nil {
            log.Fatal("Cannot add deployment metric: ", err)
        }

    err = putMetric(metricData, *ns, getEnv("PLUGIN_AWS_REGION","us-east-1"))
        if err != nil {
            log.Fatal("Cannot put metric data due to unexpected error.",err)
        }
    print("Operation completed successfuly.")
}
