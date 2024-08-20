package celplugin

import (
	"context"
	"fmt"
	"testing"

	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/apiserver/schema"
	"k8s.io/apiextensions-apiserver/pkg/apiserver/schema/cel"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
	"kcl-lang.io/kcl-go/pkg/plugin"
)

func primitiveType(typ, format string) schema.Structural {
	result := schema.Structural{
		Generic: schema.Generic{
			Type: typ,
		},
	}
	if len(format) != 0 {
		result.ValueValidation = &schema.ValueValidation{
			Format: format,
		}
	}
	return result
}

var (
	integerType = primitiveType("integer", "")
	int32Type   = primitiveType("integer", "int32")
	int64Type   = primitiveType("integer", "int64")
	numberType  = primitiveType("number", "")
	floatType   = primitiveType("number", "float")
	doubleType  = primitiveType("number", "double")
	stringType  = primitiveType("string", "")
	byteType    = primitiveType("string", "byte")
	booleanType = primitiveType("boolean", "")

	durationFormat = primitiveType("string", "duration")
	dateFormat     = primitiveType("string", "date")
	dateTimeFormat = primitiveType("string", "date-time")
)

func TestValidationExpressions(t *testing.T) {
	plugin.RegisterPlugin(plugin.Plugin{
		Name: "cel",
		MethodMap: map[string]plugin.MethodSpec{
			"add": {
				Body: func(args *plugin.MethodArgs) (*plugin.MethodResult, error) {
					v := args.IntArg(0) + args.IntArg(1)
					return &plugin.MethodResult{V: v}, nil
				},
			},
		},
	})

	s := objectTypePtr(map[string]schema.Structural{
		"presentObj": objectType(map[string]schema.Structural{
			"presentStr": stringType,
		}),
		"absentObj": objectType(map[string]schema.Structural{
			"absentStr": stringType,
		}),
		"m": mapType(&stringType),
		"l": listType(&stringType),
	})
	s.XValidations = apiextensions.ValidationRules{
		{
			Rule:   `self.m[?'k'].orValue('') == 'v2'`,
			Reason: ptr.To(apiextensions.FieldValueDuplicate),
			// FieldPath: ".field2",
		},
	}
	validator := cel.NewValidator(s, true, uint64(10000))
	errs, _ := validator.Validate(
		context.TODO(),
		field.NewPath("root"),
		s,
		map[string]interface{}{
			"presentObj": map[string]interface{}{
				"presentStr": "value",
			},
			"m": map[string]interface{}{"k": "v"},
			"l": []interface{}{"a"},
		},
		nil,
		1000)

	for _, err := range errs {
		t.Errorf("unexpected error: %v", err)
	}

}

func objs(val ...interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(val))
	for i, v := range val {
		result[fmt.Sprintf("val%d", i+1)] = v
	}
	return result
}

func schemas(valSchema ...schema.Structural) *schema.Structural {
	result := make(map[string]schema.Structural, len(valSchema))
	for i, v := range valSchema {
		result[fmt.Sprintf("val%d", i+1)] = v
	}
	return objectTypePtr(result)
}

func listType(items *schema.Structural) schema.Structural {
	return arrayType("atomic", nil, items)
}

func listTypePtr(items *schema.Structural) *schema.Structural {
	l := listType(items)
	return &l
}

func listSetType(items *schema.Structural) schema.Structural {
	return arrayType("set", nil, items)
}

func listMapType(keys []string, items *schema.Structural) schema.Structural {
	return arrayType("map", keys, items)
}

func listMapTypePtr(keys []string, items *schema.Structural) *schema.Structural {
	l := listMapType(keys, items)
	return &l
}

func arrayType(listType string, keys []string, items *schema.Structural) schema.Structural {
	result := schema.Structural{
		Generic: schema.Generic{
			Type: "array",
		},
		Extensions: schema.Extensions{
			XListType: &listType,
		},
		Items: items,
	}
	if len(keys) > 0 && listType == "map" {
		result.Extensions.XListMapKeys = keys
	}
	return result
}

func objectType(props map[string]schema.Structural) schema.Structural {
	return schema.Structural{
		Generic: schema.Generic{
			Type: "object",
		},
		Properties: props,
	}
}

func objectTypePtr(props map[string]schema.Structural) *schema.Structural {
	o := objectType(props)
	return &o
}

func mapType(valSchema *schema.Structural) schema.Structural {
	result := schema.Structural{
		Generic: schema.Generic{
			Type: "object",
		},
		AdditionalProperties: &schema.StructuralOrBool{Bool: true, Structural: valSchema},
	}
	return result
}

// func mapTypePtr(valSchema *schema.Structural) *schema.Structural {
// 	m := mapType(valSchema)
// 	return &m
// }

func intOrStringType() schema.Structural {
	return schema.Structural{
		Extensions: schema.Extensions{
			XIntOrString: true,
		},
	}
}
