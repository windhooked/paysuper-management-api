package manager

import (
	"github.com/ProtocolONE/p1pay.api/database/dao"
	"github.com/ProtocolONE/p1pay.api/database/model"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	minPaymentAmount float64 = 0
	maxPaymentAmount float64 = 15000
)

type ProjectManager struct {
	*Manager

	merchantManager *MerchantManager
	currencyManager *CurrencyManager
}

func InitProjectManager(database dao.Database, logger *zap.SugaredLogger) *ProjectManager {
	pm := &ProjectManager{
		Manager:         &Manager{Database: database, Logger: logger},
		merchantManager: InitMerchantManager(database, logger),
		currencyManager: InitCurrencyManager(database, logger),
	}

	return pm
}

func (pm *ProjectManager) Create(ps *model.ProjectScalar) (*model.Project, error) {
	p := &model.Project{
		Id:                         bson.NewObjectId(),
		Merchant:                   ps.Merchant,
		Name:                       ps.Name,
		CallbackProtocol:           ps.CallbackProtocol,
		CreateInvoiceAllowedUrls:   ps.CreateInvoiceAllowedUrls,
		IsAllowDynamicNotifyUrls:   ps.IsAllowDynamicNotifyUrls,
		IsAllowDynamicRedirectUrls: ps.IsAllowDynamicRedirectUrls,
		OnlyFixedAmounts:           ps.OnlyFixedAmounts,
		SecretKey:                  ps.SecretKey,
		URLCheckAccount:            ps.URLCheckAccount,
		URLProcessPayment:          ps.URLProcessPayment,
		URLRedirectFail:            ps.URLRedirectFail,
		URLRedirectSuccess:         ps.URLRedirectSuccess,
		SendNotifyEmail:            ps.SendNotifyEmail,
		NotifyEmails:               ps.NotifyEmails,
		IsActive:                   ps.IsActive,
		CreatedAt:                  time.Now(),
		UpdatedAt:                  time.Now(),
	}

	p.LimitsCurrency = p.Merchant.Currency
	p.CallbackCurrency = p.Merchant.Currency

	if ps.LimitsCurrency != nil {
		if c := pm.currencyManager.FindByCodeInt(*ps.LimitsCurrency); c != nil {
			p.LimitsCurrency = c
		}
	}

	if ps.CallbackCurrency != nil {
		if c := pm.currencyManager.FindByCodeInt(*ps.CallbackCurrency); c != nil {
			p.CallbackCurrency = c
		}
	}

	p.MinPaymentAmount = minPaymentAmount
	p.MaxPaymentAmount = maxPaymentAmount

	if ps.MinPaymentAmount != nil {
		p.MinPaymentAmount = *ps.MinPaymentAmount
	}

	if ps.MaxPaymentAmount != nil {
		p.MaxPaymentAmount = *ps.MaxPaymentAmount
	}

	err := pm.Database.Repository(tableProject).InsertProject(p)

	if err != nil {
		pm.Logger.Errorf("Query from table \"%s\" ended with error: %s", tableProject, err)
	}

	return p, err
}

func (pm *ProjectManager) Update(p *model.Project, pn *model.ProjectScalar) (*model.Project, error) {
	p.CreateInvoiceAllowedUrls = pn.CreateInvoiceAllowedUrls
	p.NotifyEmails = pn.NotifyEmails
	p.UpdatedAt = time.Now()

	if p.Name != pn.Name {
		p.Name = pn.Name
	}

	if p.CallbackProtocol != pn.CallbackProtocol {
		p.CallbackProtocol = pn.CallbackProtocol
	}

	if p.IsAllowDynamicNotifyUrls != pn.IsAllowDynamicNotifyUrls {
		p.IsAllowDynamicNotifyUrls = pn.IsAllowDynamicNotifyUrls
	}

	if p.IsAllowDynamicRedirectUrls != pn.IsAllowDynamicRedirectUrls {
		p.IsAllowDynamicRedirectUrls = pn.IsAllowDynamicRedirectUrls
	}

	if p.OnlyFixedAmounts != pn.OnlyFixedAmounts {
		p.OnlyFixedAmounts = pn.OnlyFixedAmounts
	}

	if p.SecretKey != pn.SecretKey {
		p.SecretKey = pn.SecretKey
	}

	if p.URLCheckAccount != pn.URLCheckAccount {
		p.URLCheckAccount = pn.URLCheckAccount
	}

	if p.URLProcessPayment != pn.URLProcessPayment {
		p.URLProcessPayment = pn.URLProcessPayment
	}

	if p.URLRedirectFail != pn.URLRedirectFail {
		p.URLRedirectFail = pn.URLRedirectFail
	}

	if p.URLRedirectSuccess != pn.URLRedirectSuccess {
		p.URLRedirectSuccess = pn.URLRedirectSuccess
	}

	if p.SendNotifyEmail != pn.SendNotifyEmail {
		p.SendNotifyEmail = pn.SendNotifyEmail
	}

	if p.IsActive != pn.IsActive {
		p.IsActive = pn.IsActive
	}

	if pn.LimitsCurrency != nil && (p.LimitsCurrency == nil || p.LimitsCurrency.CodeInt != *pn.LimitsCurrency) {
		if c := pm.currencyManager.FindByCodeInt(*pn.LimitsCurrency); c != nil {
			p.LimitsCurrency = c
		}
	}

	if pn.CallbackCurrency != nil && (p.CallbackCurrency == nil || p.CallbackCurrency.CodeInt != *pn.CallbackCurrency) {
		if c := pm.currencyManager.FindByCodeInt(*pn.CallbackCurrency); c != nil {
			p.CallbackCurrency = c
		}
	}

	if pn.MinPaymentAmount != nil && p.MinPaymentAmount != *pn.MinPaymentAmount {
		p.MinPaymentAmount = *pn.MinPaymentAmount
	}

	if pn.MaxPaymentAmount != nil && p.MaxPaymentAmount != *pn.MaxPaymentAmount {
		p.MaxPaymentAmount = *pn.MaxPaymentAmount
	}

	err := pm.Database.Repository(tableProject).UpdateProject(p)

	if err != nil {
		pm.Logger.Errorf("Query from table \"%s\" ended with error: %s", tableProject, err)
	}

	return p, err
}

func (pm *ProjectManager) Delete(p *model.Project) error {
	p.IsActive = false

	return pm.Database.Repository(tableProject).UpdateProject(p)
}

func (pm *ProjectManager) FindProjectsByMerchantIdAndName(mId bson.ObjectId, pName string) *model.Project {
	return pm.Database.Repository(tableProject).FindProjectsByMerchantIdAndName(mId, pName)
}

func (pm *ProjectManager) FindProjectsByMerchantId(mId bson.ObjectId) []*model.Project {
	p, err := pm.Database.Repository(tableProject).FindProjectsByMerchantId(mId)

	if err != nil {
		pm.Logger.Errorf("Query from table \"%s\" ended with error: %s", tableProject, err)
	}

	if p == nil {
		return []*model.Project{}
	}

	return p
}

func (pm *ProjectManager) FindProjectById(id string) *model.Project {
	p, err := pm.Database.Repository(tableCurrency).FindProjectById(bson.ObjectIdHex(id))

	if err != nil {
		pm.Logger.Errorf("Query from table \"%s\" ended with error: %s", tableCurrency, err)
	}

	return p
}
