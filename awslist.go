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
            instances = append(instances, i) //TODO: remove reallocations
        }
    }

    return instances
}

func Tags2Map(ts []*ec2.Tag) (tags map[string]string) {
    tags = map[string]string{}
    for _, t := range ts {
        tags[*t.Key] = *t.Value
    }
    return tags
}

func GetVal(ts map[string]string, n string) string {
    if v, ok := ts[n]; ok == true {
        return v
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
        tags := Tags2Map(i.Tags)
        table.Append([]string{GetVal(tags, "Name"), 
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