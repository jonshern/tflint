package evaluator

import (
	"reflect"
	"testing"

	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/parser"
	"github.com/wata727/tflint/config"
	"github.com/wata727/tflint/schema"
)

func TestIsEvaluable(t *testing.T) {
	cases := []struct {
		Name   string
		Input  string
		Result bool
	}{
		{
			Name:   "var syntax",
			Input:  "${var.text}",
			Result: true,
		},
		{
			Name:   "plain text",
			Input:  "text",
			Result: true,
		},
		{
			Name:   "module syntax",
			Input:  "${module.text}",
			Result: false,
		},
		{
			Name:   "resource syntax",
			Input:  "${aws_subnet.app.id}",
			Result: false,
		},
		{
			Name:   "function syntax",
			Input:  "${lookup(var.roles, count.index)}",
			Result: false,
		},
	}

	for _, tc := range cases {
		result := isEvaluable(tc.Input)
		if result != tc.Result {
			t.Fatalf("\nBad: %t\nExpected: %t\n\ntestcase: %s", result, tc.Result, tc.Name)
		}
	}
}

func TestEvalReturnString(t *testing.T) {
	cases := []struct {
		Name   string
		Input  string
		Src    string
		Result string
	}{
		{
			Name: "completed string variable",
			Input: `
variable "name" {
    type = "string"
    default = "test"
}`,
			Src:    "${var.name}",
			Result: "test",
		},
		{
			Name: "completed list variable",
			Input: `
variable "name" {
    type = "list"
    default = ["test1", "test2"]
}`,
			Src:    "${var.name[0]}",
			Result: "test1",
		},
		{
			Name: "completed map variable",
			Input: `
variable "name" {
    type = "map"
    default = {
        key = "test1"
        value = "test2"
    }
}`,
			Src:    "${var.name[\"key\"]}",
			Result: "test1",
		},
		{
			Name: "string variable in missing type",
			Input: `
variable "name" {
    default = "test"
}`,
			Src:    "${var.name}",
			Result: "test",
		},
		{
			Name: "list variable in missing key",
			Input: `
variable "name" {
    default = ["test1", "test2"]
}`,
			Src:    "${var.name[0]}",
			Result: "test1",
		},
		{
			Name: "map variable in missing key",
			Input: `
variable "name" {
    default = {
        key = "test1"
        value = "test2"
    }
}`,
			Src:    "${var.name[\"key\"]}",
			Result: "test1",
		},
	}

	for _, tc := range cases {
		root, _ := parser.Parse([]byte(tc.Input))
		template := map[string]*ast.File{"testfile": root}

		evaluator, err := NewEvaluator(template, []*schema.Template{}, []*ast.File{}, config.Init())
		if err != nil {
			t.Fatalf("\nError: %s\n\ntestcase: %s", err, tc.Name)
		}
		result, _ := evaluator.Eval(tc.Src)
		if result != tc.Result {
			t.Fatalf("\nBad: %s\nExpected: %s\n\ntestcase: %s", result, tc.Result, tc.Name)
		}
	}
}

func TestEvalReturnList(t *testing.T) {
	cases := []struct {
		Name   string
		Input  string
		Src    string
		Result []interface{}
	}{
		{
			Name: "return list variable",
			Input: `
variable "name" {
    default = ["test1", "test2"]
}`,
			Src:    "${var.name}",
			Result: []interface{}{"test1", "test2"},
		},
	}

	for _, tc := range cases {
		root, _ := parser.Parse([]byte(tc.Input))
		template := map[string]*ast.File{"testfile": root}

		evaluator, err := NewEvaluator(template, []*schema.Template{}, []*ast.File{}, config.Init())
		if err != nil {
			t.Fatalf("\nError: %s\n\ntestcase: %s", err, tc.Name)
		}
		result, _ := evaluator.Eval(tc.Src)
		if !reflect.DeepEqual(result, tc.Result) {
			t.Fatalf("\nBad: %s\nExpected: %s\n\ntestcase: %s", result, tc.Result, tc.Name)
		}
	}
}

func TestEvalReturnMap(t *testing.T) {
	cases := []struct {
		Name   string
		Input  string
		Src    string
		Result map[string]interface{}
	}{
		{
			Name: "return map variable",
			Input: `
variable "name" {
    default = {
        key = "test1"
        value = "test2"
    }
}`,
			Src:    "${var.name}",
			Result: map[string]interface{}{"key": "test1", "value": "test2"},
		},
	}

	for _, tc := range cases {
		root, _ := parser.Parse([]byte(tc.Input))
		template := map[string]*ast.File{"testfile": root}

		evaluator, err := NewEvaluator(template, []*schema.Template{}, []*ast.File{}, config.Init())
		if err != nil {
			t.Fatalf("\nError: %s\n\ntestcase: %s", err, tc.Name)
		}
		result, _ := evaluator.Eval(tc.Src)
		if !reflect.DeepEqual(result, tc.Result) {
			t.Fatalf("\nBad: %s\nExpected: %s\n\ntestcase: %s", result, tc.Result, tc.Name)
		}
	}
}

func TestEvalReturnNil(t *testing.T) {
	cases := []struct {
		Name  string
		Input string
		Src   string
	}{
		{
			Name:  "undefined variable",
			Input: "",
			Src:   "${var.name}",
		},
		{
			Name:  "missing default",
			Input: "variable \"name\" {}",
			Src:   "${var.name}",
		},
	}

	for _, tc := range cases {
		root, _ := parser.Parse([]byte(tc.Input))
		template := map[string]*ast.File{"testfile": root}

		evaluator, err := NewEvaluator(template, []*schema.Template{}, []*ast.File{}, config.Init())
		if err != nil {
			t.Fatalf("\nError: %s\n\ntestcase: %s", err, tc.Name)
		}
		result, _ := evaluator.Eval(tc.Src)
		if result != nil {
			t.Fatalf("\nBad: %s\nExpected: nil\n\ntestcase: %s", result, tc.Name)
		}
	}
}
