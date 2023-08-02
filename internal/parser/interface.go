package parser

import (
	"fmt"
	"go/types"

	"github.com/webrpc/webrpc/schema"
)

func (p *Parser) ParseInterfaceMethods(iface *types.Interface, name string) error {
	service := &schema.Service{
		Name:   name,
		Schema: p.Schema, // denormalize/back-reference
	}

	// Loop over the interface's methods.
	for i := 0; i < iface.NumMethods(); i++ {
		method := iface.Method(i)
		if !method.Exported() {
			continue
		}

		methodName := method.Id()

		methodSignature, ok := method.Type().(*types.Signature)
		if !ok {
			return fmt.Errorf("%v(): failed to get method signature", methodName)
		}

		methodParams := methodSignature.Params()
		inputs, err := p.getMethodArguments(methodParams, true)
		if err != nil {
			return fmt.Errorf("%v(): failed to get inputs: %w", methodName, err)
		}

		// First method argument must be of type context.Context.
		if methodParams.Len() == 0 {
			return fmt.Errorf("%v(): first method argument must be context.Context: no arguments defined", methodName)
		}
		if err := ensureContextType(methodParams.At(0).Type()); err != nil {
			return fmt.Errorf("%v(): first method argument must be context.Context: %w", methodName, err)
		}
		inputs = inputs[1:] // Cut it off. The gen/golang adds context.Context as first method argument automatically.

		methodResults := methodSignature.Results()
		outputs, err := p.getMethodArguments(methodResults, false)
		if err != nil {
			return fmt.Errorf("%v(): failed to get outputs: %w", methodName, err)
		}

		// Last method return value must be of type error.
		if methodResults.Len() == 0 {
			return fmt.Errorf("%v(): last return value must be context.Context: no return values defined", methodName)
		}
		if err := ensureErrorType(methodResults.At(methodResults.Len() - 1).Type()); err != nil {
			return fmt.Errorf("%v(): first method argument must be context.Context: %w", methodName, err)
		}
		outputs = outputs[:len(outputs)-1] // Cut it off. The gen/golang adds error as a last return value automatically.

		service.Methods = append(service.Methods, &schema.Method{
			Name:    methodName,
			Inputs:  inputs,
			Outputs: outputs,
			Service: service, // denormalize/back-reference
		})
	}

	if len(service.Methods) == 0 {
		// Ignore interfaces with no methods defined.
		return nil
	}

	p.Schema.Services = append(p.Schema.Services, service)
	return nil
}

func (p *Parser) getMethodArguments(params *types.Tuple, isInput bool) ([]*schema.MethodArgument, error) {
	var args []*schema.MethodArgument

	for i := 0; i < params.Len(); i++ {
		param := params.At(i)
		typ := param.Type()

		name := param.Name()
		if name == "" {
			// TODO: Name the field based on field type? (strings []string, hub *Hub)
			if isInput {
				name = fmt.Sprintf("arg%v", i) // 0 is `ctx context.Context`
			} else {
				name = fmt.Sprintf("ret%v", i+1)
			}
		}

		varType, err := p.ParseType(typ) // Type name will be resolved deeper down the stack.
		if err != nil {
			return nil, fmt.Errorf("failed to parse argument %v %v: %w", name, typ, err)
		}

		arg := &schema.MethodArgument{
			Name:      name,
			Type:      varType,
			InputArg:  isInput,  // denormalize/back-reference
			OutputArg: !isInput, // denormalize/back-reference
		}

		args = append(args, arg)
	}

	return args, nil
}

func ensureContextType(typ types.Type) (err error) {
	namedType, ok := typ.(*types.Named)
	if !ok {
		return fmt.Errorf("expected named type: found type %T (%+v)", typ, typ)
	}

	if _, ok := namedType.Underlying().(*types.Interface); !ok {
		return fmt.Errorf("expected underlying interface: found type %T (%+v)", typ, typ)
	}

	pkgName := namedType.Obj().Pkg().Name()
	typeName := namedType.Obj().Name()

	if pkgName != "context" && typeName != "Context" {
		return fmt.Errorf("expected context.Context: found %v.%v", pkgName, typeName)
	}

	return nil
}

func ensureErrorType(typ types.Type) (err error) {
	namedType, ok := typ.(*types.Named)
	if !ok {
		return fmt.Errorf("expected named type: found type %T (%+v)", typ, typ)
	}

	if _, ok := namedType.Underlying().(*types.Interface); !ok {
		return fmt.Errorf("expected underlying interface: found type %T (%+v)", typ, typ)
	}

	pkgName := namedType.Obj().Pkg()
	typeName := namedType.Obj().Name()

	if pkgName != nil && typeName != "error" {
		return fmt.Errorf("expected error: found %v.%v", pkgName, typeName)
	}

	return nil
}
