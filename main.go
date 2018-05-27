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
    var ns = flag.String("TestSpace", "Linux/System", "CloudWatch metric namespace") 
    var err = putMetric(metricData, *ns, "us-east-1")

    dim := cloudwatch.Dimension{
            Name:  aws.String(serviceName),
            Environment: aws.String(serviceEnvironment),
        }
    //Metric data addition
    metricData, err = addMetric("Deployment", "Count", deployment, dims, metricData)
        if err != nil {
            log.Fatal("Can't add memory available metric: ", err)
        }

    if err != nil {
        log.Fatal("Can't put CloudWatch Metric",err)
    }
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
