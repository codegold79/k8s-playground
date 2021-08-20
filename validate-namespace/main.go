package main

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/validation"
)

func main() {
	namespaces := []string{
		"forward/slash",
		"normal",
		"123number",
		"number123",
		"12345",
		"dashes-in-name",
		"under_score",
		"/",
		"slash/",
	}

	printNamespaceValidity(namespaces)
}

func printNamespaceValidity(namespaces []string) {
	for _, ns := range namespaces {
		if errs := validation.ValidateNamespaceName(ns, false); errs != nil {
			fmt.Println("\nerrors in", ns)
			for _, e := range errs {
				fmt.Println(e)
			}
		}
	}
}
