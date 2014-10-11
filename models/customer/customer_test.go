package customer

import (
	"testing"
	"time"
)

func TestNewCustomer(t *testing.T) {
	c := NewCustomer()

	if c == nil {
		t.Errorf("expected customer not to be nil, got %v", c)
	}
	if c.Blocked {
		t.Errorf("expected Blocked to be false, got %v", c.Blocked)
	}
	if c.RedOneShipped {
		t.Errorf("expected RedOneShipped to be false, got %v", c.RedOneShipped)
	}
	if c.EarlyAdopter {
		t.Errorf("expected EarlyAdopter to be false, got %v", c.EarlyAdopter)
	}
	if c.Multiple {
		t.Errorf("expected Multiple to be false, got %v", c.Multiple)
	}
}

func TestCustomer_IsEducationDiscountApproved(t *testing.T) {
	c := NewCustomer()
	c.ErpCode = "ESN"

	if res := c.IsEducationDiscountApproved(); !res {
		t.Errorf("expected IsEducationDiscountApproved to return true got %v", res)
	}

	c.ErpCode = "NOPE"
	if res := c.IsEducationDiscountApproved(); res {
		t.Errorf("expected IsEducationDiscountApproved to return false got %v", res)
	}
}

func TestCustomer_IsEducationDiscountDenied(t *testing.T) {
	c := NewCustomer()
	c.ErpCode = "ESV"

	if res := c.IsEducationDiscountDenied(); !res {
		t.Errorf("expected IsEducationDiscountDenied to return true, got %v", res)
	}

	c.ErpCode = "NOPE"
	if res := c.IsEducationDiscountDenied(); res {
		t.Errorf("expected IsEducationDiscountDenied to return false, got %v", res)
	}
}

func TestCustomer_IsEducationDiscountVerified(t *testing.T) {
	c := NewCustomer()
	c.ErpCode = "ETV"

	if res := c.IsEducationDiscountVerified(); !res {
		t.Errorf("expected IsEducationDiscountVerfiied to return true, got %v", res)
	}

	c.ErpCode = "NOPE"
	if res := c.IsEducationDiscountVerified(); res {
		t.Errorf("expected IsEducationDiscountVerfiied to return false, got %v", res)
	}
}

func TestCustomer_HasEducationProgram(t *testing.T) {
	c := NewCustomer()

	c.ErpCode = "ESV"
	if res := c.HasEducationProgram(); !res {
		t.Errorf("expected HasEducationProgram to return true, got %v", res)
	}

	c.ErpCode = "ESN"
	if res := c.HasEducationProgram(); !res {
		t.Errorf("expected HasEducationProgram to return true, got %v", res)
	}

	c.ErpCode = "NOPE"
	if res := c.HasEducationProgram(); res {
		t.Errorf("expected HasEducationProgram to return false, got %v", res)
	}
}

func TestCustomer_isEducationDiscountExpired(t *testing.T) {
	currentTime := time.Now()
	expiredTime := currentTime.AddDate(-1, 0, 0)
	unexpiredTime := currentTime.AddDate(1, 0, 0)
	c := NewCustomer()

	c.EducationExpiresAt = expiredTime
	if res := c.isEducationDiscountExpired(); !res {
		t.Errorf("expected isEducationDiscountExpired to be true, got %v", res)
	}

	c.EducationExpiresAt = unexpiredTime
	if res := c.isEducationDiscountExpired(); res {
		t.Errorf("expected isEducationDiscountExpired to be false, got %v", res)
	}
}

func TestCustomer_isInDiscountGroup(t *testing.T) {
	currentTime := time.Now()
	expiredTime := currentTime.AddDate(-1, 0, 0)
	unexpiredTime := currentTime.AddDate(1, 0, 0)
	badDiscountCode := "NOPE"
	goodDiscountCode := "ESV"
	c := NewCustomer()

	c.ErpCode = badDiscountCode
	if res := c.isInDiscountGroup(); res {
		t.Errorf("expected isInDiscountGroup to be false, got %v", res)
	}

	c.ErpCode = goodDiscountCode
	c.EducationExpiresAt = expiredTime
	if res := c.isInDiscountGroup(); res {
		t.Errorf("expected isInDiscountGroup to be false, got %v", res)
	}

	c.EducationExpiresAt = unexpiredTime
	if res := c.isInDiscountGroup(); !res {
		t.Errorf("expected isInDiscountGroup to be true, got %v", res)
	}
}

func TestFindByID(t *testing.T) {
	DBSession, _ = NewSession()
	defer DBSession.Close()
	c, err := FindByID("53b2f98138fc31015600000a")

	if c == nil {
		t.Errorf("expected non nil customer: %s\n", err)
	}
}

// func TestCustomer_EducationDiscountStatus(t testing.T) {
//
// }
