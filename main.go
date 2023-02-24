package main

import (
	"context"
	"fmt"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vim25/soap"
	trace "govmomi-tracing/pkg/trace"
	"net/url"
	"os"
	"time"
)

func main() {

	username := os.Getenv("GOVMOMI_USERNAME")
	password := os.Getenv("GOVMOMI_PASSWORD")
	server := os.Getenv("GOVMOMI_SERVER")

	//op := trace.FromContext(ctx, "VMPathNameAsURL")

	topctx, cancel := context.WithTimeout(context.TODO(), 300*time.Second)
	defer cancel()
	op := trace.NewOperation(topctx, "logon")

	u, err := soap.ParseURL(server)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	u.User = url.UserPassword(username, password)
	c, err := govmomi.NewClient(op.Context, u, true)

	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}

	op = trace.NewOperation(topctx, "rest logon")
	restClient := rest.NewClient(c.Client)
	err = restClient.Login(op.Context, u.User)
	if err != nil {
		logoutErr := c.Logout(context.TODO())
		if logoutErr != nil {
			err = logoutErr
		}
		fmt.Fprint(os.Stderr, err)
	}

	for {
		var datacenterPath string
		finder := find.NewFinder(c.Client, true)

		op = trace.NewOperation(topctx, "DatacenterList")
		datacenters, err := finder.DatacenterList(op.Context, "./...")
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		}

		for _, dc := range datacenters {
			fmt.Println(dc.Name())
			fmt.Println(dc.InventoryPath)

			datacenterPath = dc.InventoryPath
		}

		fmt.Println("Sleeping for 30 seconds...")
		time.Sleep(30 * time.Second)

		op = trace.NewOperation(topctx, "ClusterComputeResourceList")
		clusters, err := finder.ClusterComputeResourceList(op.Context, "./...")
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		}

		for _, c := range clusters {
			fmt.Println(c.Name())
			fmt.Println(c.InventoryPath)
		}

		fmt.Println("Sleeping for 30 seconds...")
		time.Sleep(30 * time.Second)

		op = trace.NewOperation(topctx, "VirtualMachineList")
		virtualMachines, err := finder.VirtualMachineList(op.Context, datacenterPath+"/vm/")
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		}

		for _, vm := range virtualMachines {
			fmt.Println(vm.Name())
			fmt.Println(vm.InventoryPath)
		}
		fmt.Println("Sleeping for 300 seconds...")
		time.Sleep(300 * time.Second)
	}
}
