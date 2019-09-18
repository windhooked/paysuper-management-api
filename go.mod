module github.com/paysuper/paysuper-management-api

require (
	github.com/Jeffail/gabs v1.1.1 // indirect
	github.com/ProtocolONE/authone-jwt-verifier-golang v0.0.0-20190327070329-4dd563b01681
	github.com/ProtocolONE/geoip-service v0.0.0-20190903084234-1d5ae6b96679
	github.com/SAP/go-hdb v0.13.2 // indirect
	github.com/SebastiaanKlippert/go-wkhtmltopdf v1.4.1
	github.com/SermoDigital/jose v0.9.2-0.20161205224733-f6df55f235c2 // indirect
	github.com/asaskevich/govalidator v0.0.0-20180720115003-f9ffefc3facf // indirect
	github.com/aws/aws-sdk-go v1.23.16
	github.com/elazarl/go-bindata-assetfs v1.0.0 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/golang/protobuf v1.3.2
	github.com/google/uuid v1.1.1
	github.com/hashicorp/go-memdb v0.0.0-20181108192425-032f93b25bec // indirect
	github.com/hashicorp/go-plugin v0.0.0-20181212150838-f444068e8f5a // indirect
	github.com/karlseguin/expect v1.0.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/keybase/go-crypto v0.0.0-20181127160227-255a5089e85a // indirect
	github.com/labstack/echo/v4 v4.1.6
	github.com/micro/go-micro v1.8.0
	github.com/micro/go-plugins v1.2.0
	github.com/micro/go-rcache v0.2.1 // indirect
	github.com/micro/util v0.2.0 // indirect
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/paysuper/echo-casbin-middleware v0.0.0-20190918141650-c6a2b0e617b8
	github.com/paysuper/paysuper-aws-manager v0.0.1
	github.com/paysuper/paysuper-billing-server v0.0.0-20190918083056-4fbf54925f32
	github.com/paysuper/paysuper-payment-link v0.0.0-20190903143854-b799a77c03ce
	github.com/paysuper/paysuper-recurring-repository v1.0.123
	github.com/paysuper/paysuper-reporter v0.0.0-20190917180039-6701d139ca7f
	github.com/paysuper/paysuper-tax-service v0.0.0-20190903084038-7849f394f122
	github.com/stretchr/testify v1.4.0
	github.com/ttacon/builder v0.0.0-20170518171403-c099f663e1c2 // indirect
	github.com/ttacon/libphonenumber v1.0.1
	github.com/wsxiaoys/terminal v0.0.0-20160513160801-0940f3fc43a0 // indirect
	go.uber.org/zap v1.10.0
	gopkg.in/go-playground/validator.v9 v9.29.1
	gopkg.in/karlseguin/expect.v1 v1.0.1 // indirect
)

replace (
	github.com/gogo/protobuf v0.0.0-20190410021324-65acae22fc9 => github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d
	github.com/hashicorp/consul => github.com/hashicorp/consul v1.5.1
	github.com/lucas-clemente/quic-go => github.com/lucas-clemente/quic-go v0.12.0
	github.com/marten-seemann/qtls => github.com/marten-seemann/qtls v0.3.2
	gopkg.in/DATA-DOG/go-sqlmock.v1 => github.com/DATA-DOG/go-sqlmock v1.3.3
)

go 1.12
