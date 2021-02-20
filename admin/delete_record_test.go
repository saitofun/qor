package admin_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	. "github.com/saitofun/qor/admin/tests/dummy"
	"github.com/saitofun/qor/gorm"
)

func TestDeleteRecord(t *testing.T) {
	user := User{Name: "delete_record", Role: "admin"}
	db.Save(&user)
	form := url.Values{
		"_method": {"delete"},
	}

	if req, err := http.PostForm(server.URL+"/admin/users/"+fmt.Sprint(user.ID), form); err == nil {
		if req.StatusCode != 200 {
			t.Errorf("Delete request should be processed successfully")
		}

		if !errors.Is(db.First(&User{}, "name = ?", "delete_record").Error, gorm.ErrRecordNotFound) {
			t.Errorf("User should be deleted successfully")
		}
	} else {
		t.Errorf(err.Error())
	}
}
