package putio

import (
	"fmt"
	"net/http"
	"testing"
)

func TestAccount_Info(t *testing.T) {
	setup()
	defer teardown()

	fixture := `
{
	"info": {
		"account_active": true,
		"avatar_url": "https://some-valid-gravatar-url.com/avatar.jpg?s=50",
		"days_until_files_deletion": 0,
		"default_subtitle_language": "tur",
		"disk": {
			"avail": 715975016136,
			"size": 2147483648000,
			"used": 1431508631864
		},
		"has_voucher": 0,
		"mail": "naber@iyidir.com",
		"passive_account": false,
		"plan_expiration_date": "2020-01-01T00:00:00",
		"settings": {
			"callback_url": "",
			"default_download_folder": 0,
			"default_subtitle_language": "tur",
			"download_folder_unset": false,
			"is_invisible": false,
			"nextepisode": true,
			"private_download_host_ip": null,
			"pushover_token": "",
			"routing": "Istanbul",
			"sorting": "NAME_ASC",
			"ssl_enabled": true,
			"start_from": true,
			"subtitle_languages": [
				"tur",
				"eng"
				],
			"trSorting": "NAME_ASC",
			"use_soon": true
		},
		"simultaneous_download_limit": 100,
		"subtitle_languages": [
			"tur",
			"eng"
		],
		"user_id": 1,
		"username": "naber"
	},
	"status": "OK"
}
`
	mux.HandleFunc("/v2/account/info", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintln(w, fixture)
	})

	info, err := client.Account.Info(nil)
	if err != nil {
		t.Error(err)
	}

	if info.Username != "naber" {
		t.Errorf("got: %v, want: naber", info.Username)
	}

	if info.Mail != "naber@iyidir.com" {
		t.Errorf("got: %v, want: naber@iyidir.com", info.Mail)
	}
}

func TestAccount_Settings(t *testing.T) {
	setup()
	defer teardown()

	fixture := `
{
	"settings": {
		"callback_url": "",
		"default_download_folder": 666,
		"default_subtitle_language": "tur",
		"download_folder_unset": false,
		"is_invisible": false,
		"nextepisode": true,
		"private_download_host_ip": null,
		"pushover_token": "",
		"routing": "Istanbul",
		"sorting": "NAME_ASC",
		"ssl_enabled": true,
		"start_from": true,
		"subtitle_languages": [
			"tur",
			"eng"
		],
		"trSorting": "NAME_ASC",
		"use_soon": true
	},
	"status": "OK"
}
`
	mux.HandleFunc("/v2/account/settings", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintln(w, fixture)
	})

	settings, err := client.Account.Settings(nil)
	if err != nil {
		t.Error(err)
	}

	if settings.DefaultDownloadFolder != 666 {
		t.Errorf("got: %v, want: 0", settings.DefaultDownloadFolder)
	}

	if settings.DefaultSubtitleLanguage != "tur" {
		t.Errorf("got: %v, want: tr", settings.DefaultSubtitleLanguage)
	}
}
