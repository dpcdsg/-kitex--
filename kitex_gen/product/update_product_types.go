// Manual additions aligned with idl/product.thrift (UpdateProduct). Regenerate with kitex when possible.

package product

import (
	"fmt"
	thrift "github.com/cloudwego/kitex/pkg/protocol/bthrift/apache"
)

type UpdateProductRequest struct {
	Token       string `thrift:"token,1" frugal:"1,default,string" json:"token"`
	ProductId   int64  `thrift:"product_id,2" frugal:"2,default,i64" json:"product_id"`
	Name        string `thrift:"name,3" frugal:"3,default,string" json:"name"`
	Description string `thrift:"description,4" frugal:"4,default,string" json:"description"`
	Price       int64  `thrift:"price,5" frugal:"5,default,i64" json:"price"`
	Stock       int64  `thrift:"stock,6" frugal:"6,default,i64" json:"stock"`
	ImageUrl    string `thrift:"image_url,7" frugal:"7,default,string" json:"image_url"`
	Category    string `thrift:"category,8" frugal:"8,default,string" json:"category"`
}

func NewUpdateProductRequest() *UpdateProductRequest { return &UpdateProductRequest{} }

func (p *UpdateProductRequest) InitDefault() {}

func (p *UpdateProductRequest) GetToken() string       { return p.Token }
func (p *UpdateProductRequest) GetProductId() int64    { return p.ProductId }
func (p *UpdateProductRequest) GetName() string        { return p.Name }
func (p *UpdateProductRequest) GetDescription() string { return p.Description }
func (p *UpdateProductRequest) GetPrice() int64        { return p.Price }
func (p *UpdateProductRequest) GetStock() int64        { return p.Stock }
func (p *UpdateProductRequest) GetImageUrl() string    { return p.ImageUrl }
func (p *UpdateProductRequest) GetCategory() string    { return p.Category }

func (p *UpdateProductRequest) SetToken(v string)        { p.Token = v }
func (p *UpdateProductRequest) SetProductId(v int64)     { p.ProductId = v }
func (p *UpdateProductRequest) SetName(v string)         { p.Name = v }
func (p *UpdateProductRequest) SetDescription(v string)  { p.Description = v }
func (p *UpdateProductRequest) SetPrice(v int64)         { p.Price = v }
func (p *UpdateProductRequest) SetStock(v int64)         { p.Stock = v }
func (p *UpdateProductRequest) SetImageUrl(v string)     { p.ImageUrl = v }
func (p *UpdateProductRequest) SetCategory(v string)     { p.Category = v }

var fieldIDToName_UpdateProductRequest = map[int16]string{
	1: "token", 2: "product_id", 3: "name", 4: "description",
	5: "price", 6: "stock", 7: "image_url", 8: "category",
}

func (p *UpdateProductRequest) Read(iprot thrift.TProtocol) (err error) {
	var fieldTypeId thrift.TType
	var fieldId int16
	if _, err = iprot.ReadStructBegin(); err != nil {
		goto ReadStructBeginError
	}
	for {
		_, fieldTypeId, fieldId, err = iprot.ReadFieldBegin()
		if err != nil {
			goto ReadFieldBeginError
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if fieldTypeId == thrift.STRING {
				if err = p.ReadField1(iprot); err != nil {
					goto ReadFieldError
				}
			} else if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		case 2:
			if fieldTypeId == thrift.I64 {
				if err = p.ReadField2(iprot); err != nil {
					goto ReadFieldError
				}
			} else if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		case 3:
			if fieldTypeId == thrift.STRING {
				if err = p.ReadField3(iprot); err != nil {
					goto ReadFieldError
				}
			} else if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		case 4:
			if fieldTypeId == thrift.STRING {
				if err = p.ReadField4(iprot); err != nil {
					goto ReadFieldError
				}
			} else if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		case 5:
			if fieldTypeId == thrift.I64 {
				if err = p.ReadField5(iprot); err != nil {
					goto ReadFieldError
				}
			} else if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		case 6:
			if fieldTypeId == thrift.I64 {
				if err = p.ReadField6(iprot); err != nil {
					goto ReadFieldError
				}
			} else if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		case 7:
			if fieldTypeId == thrift.STRING {
				if err = p.ReadField7(iprot); err != nil {
					goto ReadFieldError
				}
			} else if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		case 8:
			if fieldTypeId == thrift.STRING {
				if err = p.ReadField8(iprot); err != nil {
					goto ReadFieldError
				}
			} else if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		default:
			if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		}
		if err = iprot.ReadFieldEnd(); err != nil {
			goto ReadFieldEndError
		}
	}
	if err = iprot.ReadStructEnd(); err != nil {
		goto ReadStructEndError
	}
	return nil
ReadStructBeginError:
	return thrift.PrependError(fmt.Sprintf("%T read struct begin error: ", p), err)
ReadFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T read field %d begin error: ", p, fieldId), err)
ReadFieldError:
	return thrift.PrependError(fmt.Sprintf("%T read field %d '%s' error: ", p, fieldId, fieldIDToName_UpdateProductRequest[fieldId]), err)
SkipFieldError:
	return thrift.PrependError(fmt.Sprintf("%T field %d skip type %d error: ", p, fieldId, fieldTypeId), err)
ReadFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T read field end error", p), err)
ReadStructEndError:
	return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
}

func (p *UpdateProductRequest) ReadField1(iprot thrift.TProtocol) error {
	v, err := iprot.ReadString()
	if err != nil {
		return err
	}
	p.Token = v
	return nil
}
func (p *UpdateProductRequest) ReadField2(iprot thrift.TProtocol) error {
	v, err := iprot.ReadI64()
	if err != nil {
		return err
	}
	p.ProductId = v
	return nil
}
func (p *UpdateProductRequest) ReadField3(iprot thrift.TProtocol) error {
	v, err := iprot.ReadString()
	if err != nil {
		return err
	}
	p.Name = v
	return nil
}
func (p *UpdateProductRequest) ReadField4(iprot thrift.TProtocol) error {
	v, err := iprot.ReadString()
	if err != nil {
		return err
	}
	p.Description = v
	return nil
}
func (p *UpdateProductRequest) ReadField5(iprot thrift.TProtocol) error {
	v, err := iprot.ReadI64()
	if err != nil {
		return err
	}
	p.Price = v
	return nil
}
func (p *UpdateProductRequest) ReadField6(iprot thrift.TProtocol) error {
	v, err := iprot.ReadI64()
	if err != nil {
		return err
	}
	p.Stock = v
	return nil
}
func (p *UpdateProductRequest) ReadField7(iprot thrift.TProtocol) error {
	v, err := iprot.ReadString()
	if err != nil {
		return err
	}
	p.ImageUrl = v
	return nil
}
func (p *UpdateProductRequest) ReadField8(iprot thrift.TProtocol) error {
	v, err := iprot.ReadString()
	if err != nil {
		return err
	}
	p.Category = v
	return nil
}

func (p *UpdateProductRequest) Write(oprot thrift.TProtocol) (err error) {
	var fieldId int16
	if err = oprot.WriteStructBegin("UpdateProductRequest"); err != nil {
		goto WriteStructBeginError
	}
	if p != nil {
		if err = p.writeField1(oprot); err != nil {
			fieldId = 1
			goto WriteFieldError
		}
		if err = p.writeField2(oprot); err != nil {
			fieldId = 2
			goto WriteFieldError
		}
		if err = p.writeField3(oprot); err != nil {
			fieldId = 3
			goto WriteFieldError
		}
		if err = p.writeField4(oprot); err != nil {
			fieldId = 4
			goto WriteFieldError
		}
		if err = p.writeField5(oprot); err != nil {
			fieldId = 5
			goto WriteFieldError
		}
		if err = p.writeField6(oprot); err != nil {
			fieldId = 6
			goto WriteFieldError
		}
		if err = p.writeField7(oprot); err != nil {
			fieldId = 7
			goto WriteFieldError
		}
		if err = p.writeField8(oprot); err != nil {
			fieldId = 8
			goto WriteFieldError
		}
	}
	if err = oprot.WriteFieldStop(); err != nil {
		goto WriteFieldStopError
	}
	if err = oprot.WriteStructEnd(); err != nil {
		goto WriteStructEndError
	}
	return nil
WriteStructBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
WriteFieldError:
	return thrift.PrependError(fmt.Sprintf("%T write field %d error: ", p, fieldId), err)
WriteFieldStopError:
	return thrift.PrependError(fmt.Sprintf("%T write field stop error: ", p), err)
WriteStructEndError:
	return thrift.PrependError(fmt.Sprintf("%T write struct end error: ", p), err)
}

func (p *UpdateProductRequest) writeField1(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("token", thrift.STRING, 1); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteString(p.Token); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 1 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 1 end error: ", p), err)
}

func (p *UpdateProductRequest) writeField2(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("product_id", thrift.I64, 2); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteI64(p.ProductId); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 2 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 2 end error: ", p), err)
}

func (p *UpdateProductRequest) writeField3(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("name", thrift.STRING, 3); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteString(p.Name); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 3 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 3 end error: ", p), err)
}

func (p *UpdateProductRequest) writeField4(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("description", thrift.STRING, 4); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteString(p.Description); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 4 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 4 end error: ", p), err)
}

func (p *UpdateProductRequest) writeField5(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("price", thrift.I64, 5); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteI64(p.Price); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 5 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 5 end error: ", p), err)
}

func (p *UpdateProductRequest) writeField6(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("stock", thrift.I64, 6); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteI64(p.Stock); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 6 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 6 end error: ", p), err)
}

func (p *UpdateProductRequest) writeField7(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("image_url", thrift.STRING, 7); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteString(p.ImageUrl); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 7 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 7 end error: ", p), err)
}

func (p *UpdateProductRequest) writeField8(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("category", thrift.STRING, 8); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteString(p.Category); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 8 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 8 end error: ", p), err)
}

func (p *UpdateProductRequest) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("UpdateProductRequest(%+v)", *p)
}

// UpdateProductResponse matches CreateProductResponse on wire (base + product + seller_id).
type UpdateProductResponse = CreateProductResponse

type ProductServiceUpdateProductArgs struct {
	Req *UpdateProductRequest `thrift:"req,1" frugal:"1,default,UpdateProductRequest" json:"req"`
}

func NewProductServiceUpdateProductArgs() *ProductServiceUpdateProductArgs {
	return &ProductServiceUpdateProductArgs{}
}

func (p *ProductServiceUpdateProductArgs) InitDefault() {}

var ProductServiceUpdateProductArgs_Req_DEFAULT *UpdateProductRequest

func (p *ProductServiceUpdateProductArgs) GetReq() (v *UpdateProductRequest) {
	if !p.IsSetReq() {
		return ProductServiceUpdateProductArgs_Req_DEFAULT
	}
	return p.Req
}
func (p *ProductServiceUpdateProductArgs) SetReq(val *UpdateProductRequest) { p.Req = val }

var fieldIDToName_ProductServiceUpdateProductArgs = map[int16]string{1: "req"}

func (p *ProductServiceUpdateProductArgs) IsSetReq() bool { return p.Req != nil }

func (p *ProductServiceUpdateProductArgs) Read(iprot thrift.TProtocol) (err error) {
	var fieldTypeId thrift.TType
	var fieldId int16
	if _, err = iprot.ReadStructBegin(); err != nil {
		goto ReadStructBeginError
	}
	for {
		_, fieldTypeId, fieldId, err = iprot.ReadFieldBegin()
		if err != nil {
			goto ReadFieldBeginError
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if fieldTypeId == thrift.STRUCT {
				if err = p.ReadField1(iprot); err != nil {
					goto ReadFieldError
				}
			} else if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		default:
			if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		}
		if err = iprot.ReadFieldEnd(); err != nil {
			goto ReadFieldEndError
		}
	}
	if err = iprot.ReadStructEnd(); err != nil {
		goto ReadStructEndError
	}
	return nil
ReadStructBeginError:
	return thrift.PrependError(fmt.Sprintf("%T read struct begin error: ", p), err)
ReadFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T read field %d begin error: ", p, fieldId), err)
ReadFieldError:
	return thrift.PrependError(fmt.Sprintf("%T read field %d '%s' error: ", p, fieldId, fieldIDToName_ProductServiceUpdateProductArgs[fieldId]), err)
SkipFieldError:
	return thrift.PrependError(fmt.Sprintf("%T field %d skip type %d error: ", p, fieldId, fieldTypeId), err)
ReadFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T read field end error", p), err)
ReadStructEndError:
	return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
}

func (p *ProductServiceUpdateProductArgs) ReadField1(iprot thrift.TProtocol) error {
	_field := NewUpdateProductRequest()
	if err := _field.Read(iprot); err != nil {
		return err
	}
	p.Req = _field
	return nil
}

func (p *ProductServiceUpdateProductArgs) Write(oprot thrift.TProtocol) (err error) {
	var fieldId int16
	if err = oprot.WriteStructBegin("UpdateProduct_args"); err != nil {
		goto WriteStructBeginError
	}
	if p != nil {
		if err = p.writeField1(oprot); err != nil {
			fieldId = 1
			goto WriteFieldError
		}
	}
	if err = oprot.WriteFieldStop(); err != nil {
		goto WriteFieldStopError
	}
	if err = oprot.WriteStructEnd(); err != nil {
		goto WriteStructEndError
	}
	return nil
WriteStructBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
WriteFieldError:
	return thrift.PrependError(fmt.Sprintf("%T write field %d error: ", p, fieldId), err)
WriteFieldStopError:
	return thrift.PrependError(fmt.Sprintf("%T write field stop error: ", p), err)
WriteStructEndError:
	return thrift.PrependError(fmt.Sprintf("%T write struct end error: ", p), err)
}

func (p *ProductServiceUpdateProductArgs) writeField1(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("req", thrift.STRUCT, 1); err != nil {
		goto WriteFieldBeginError
	}
	if err := p.Req.Write(oprot); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 1 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 1 end error: ", p), err)
}

func (p *ProductServiceUpdateProductArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ProductServiceUpdateProductArgs(%+v)", *p)
}

type ProductServiceUpdateProductResult = ProductServiceCreateProductResult
