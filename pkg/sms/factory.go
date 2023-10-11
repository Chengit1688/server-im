package sms

type Factory struct {
}

func NewSmsFactory() *Factory {
	return new(Factory)
}
func (f *Factory) GetSms(t string) ISms {
	if t == "bao" {
		return NewBao()
	}

	return nil
}
