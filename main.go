package main

import (
	"context"
	"fmt"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vim25/soap"
	"net/url"
	"os"
	"time"
)

func main() {

	username := os.Getenv("GOVMOMI_USERNAME")
	password := os.Getenv("GOVMOMI_PASSWORD")
	server := os.Getenv("GOVMOMI_SERVER")

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	u, err := soap.ParseURL(server)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	u.User = url.UserPassword(username, password)
	c, err := govmomi.NewClient(ctx, u, true)

	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}

	restClient := rest.NewClient(c.Client)
	err = restClient.Login(ctx, u.User)
	if err != nil {
		logoutErr := c.Logout(context.TODO())
		if logoutErr != nil {
			err = logoutErr
		}
		fmt.Fprint(os.Stderr, err)
	}

	for {
		finder := find.NewFinder(c.Client, true)

		datacenters, err := finder.DatacenterList(ctx, "./...")
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		}

		for _, dc := range datacenters {
			fmt.Println(dc.Name())
			fmt.Println(dc.InventoryPath)
		}

		time.Sleep(30 * time.Second)

		clusters, err := finder.ClusterComputeResourceList(ctx, "./...")

		for _, c := range clusters {
			fmt.Println(c.Name())
			fmt.Println(c.InventoryPath)
		}
		time.Sleep(30 * time.Second)

		virtualMachines, err := finder.VirtualMachineList(ctx, "./...")
		for _, vm := range virtualMachines {
			fmt.Println(vm.Name())
			fmt.Println(vm.InventoryPath)
		}
		time.Sleep(300 * time.Second)
	}
}
