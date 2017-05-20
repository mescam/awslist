package main

import (
    "os"
    "github.com/olekukonko/tablewriter"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ec2"
)

func check(err error) {
    if err != nil {
        panic(err)
    }
}

func InstancesList(svc *ec2.EC2) (instances []*ec2.Instance) {
    resp, err := svc.DescribeInstances(nil)
    check(err)

    for _, r := range resp.Reservations {
        for _, i := range r.Instances {
            instances = append(instances, i) //TODO: remove reallactions
        }
    }

    return instances
}

func GetTagValue(ts []*ec2.Tag, name string) string {
    for _, t := range ts {
        if *t.Key == name {
            return *t.Value
        }
    }

    return ""
} 

func GenerateTable(data []*ec2.Instance) *tablewriter.Table {
    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader([]string{"Name", 
                             "Type", 
                             "State", 
                             "IP Address", 
                             "Availability Zone"})

    for _, i := range data {
        table.Append([]string{GetTagValue(i.Tags, "Name"), 
                              *i.InstanceType, 
                              *i.State.Name, 
                              *i.PrivateIpAddress, 
                              *i.Placement.AvailabilityZone})
    }

    return table
}

func main() {
    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))
    svc := ec2.New(sess)

    list := InstancesList(svc)
    GenerateTable(list).Render()
}