package main
import (
    "log"
    "fmt"
    "flag"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/awserr"
    "github.com/aws/aws-sdk-go/aws/awsutil"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/cloudwatch"
)

func main() {
    //Variable setting  
    var metricData []*cloudwatch.MetricDatum
    var ns = flag.String("namespace", "TestSpace", "CloudWatch metric namespace")//TODO:Change namespace to something meaningful.
    var dim []*cloudwatch.Dimension
    var err error

    // Hardcoded 
    var deployment float64 = 1 //TODO:Remove after functionality test"
    var serviceName = "Service name"
    var serviceNameValue = "stg-service-test"
    var serviceEnvironment = "Environment"
    var serviceEnvironmentValue = "staging"

    //Dimensions of the metrica - Service name and environment
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

    err = putMetric(metricData, *ns, "us-east-1") //TODO:Get region from credentials.
        if err != nil {
            log.Fatal("Can't put CloudWatch Metric",err)
        }
    print("Successfully added")
}

func addMetric(name, unit string, value float64, dimensions []*cloudwatch.Dimension, metricData []*cloudwatch.MetricDatum) (ret []*cloudwatch.MetricDatum, err error) {
    _metric := cloudwatch.MetricDatum{
        MetricName: aws.String(name),
        Unit:       aws.String(unit),
        Value:      aws.Float64(value),
        Dimensions: dimensions,
    }
    print("metricname",name,"unit",unit,"value",value,"dimensions",dimensions)
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
