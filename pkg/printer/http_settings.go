package printer

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const urlHttpCertServerSettings = "net/net/certificate/http.html"

var (
	errCurrentCertIdNotFound = errors.New("printer: get: failed to find current cert id")
)

// getHttpSettings fetches the HTTP Server Settings page
func (p *printer) getHttpSettings() ([]byte, error) {
	// get url & set path
	u, err := url.ParseRequestURI(p.baseUrl)
	if err != nil {
		return nil, err
	}
	u.Path = urlHttpCertServerSettings

	// make and do request
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// read body of response
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// OK status?
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("printer: get of http settings page failed (status code %d)", resp.StatusCode)
	}

	return bodyBytes, nil
}

// SetActiveCert sets the printers active certificate the specified ID and
// then restarts the printer (to make the new cert active)
// Note: This function even works of the `id` is not in the dropdown box of the printer's
// cert picker (which happens when the cert does not have a Common Name)
func (p *printer) SetActiveCert(id string) error {
	// GET http settings
	bodyBytes, err := p.getHttpSettings()
	if err != nil {
		return err
	}

	// find CSRFToken
	csrfToken, err := parseBodyForCSRFToken(bodyBytes)
	if err != nil {
		return err
	}

	// submit initial form to change the cert (MFC-L2750DW field names)
	data := url.Values{}
	data.Set("pageid", "326")
	data.Set("CSRFToken", csrfToken)
	// certificate dropdown (<select name="Bb23">)
	data.Set("Bb23", id)
	// Preserve ALL currently-enabled protocols. An absent checkbox is treated
	// as "off" by the firmware, so only flipping the HTTPS box (as upstream
	// does) would risk disabling HTTP/IPP/WebServices. The L2750DW ships these
	// enabled; we re-assert them so the only change is the certificate.
	data.Set("Ba8c", "1")        // Web Based Management HTTPS (443)
	data.Set("Ba8d", "1")        // Web Based Management HTTP  (80)
	data.Set("Ba9e", "1")        // IPP HTTPS (443)
	data.Set("ipp_ssl_used", "") // IPP secure helper (submitted empty)
	data.Set("Ba9f", "1")        // IPP HTTP (80)
	data.Set("Baa0", "1")        // IPP HTTP (631)
	data.Set("Ba7d", "1")        // Web Services HTTP
	data.Set("Bb20", "")
	data.Set("Bb21", "")
	data.Set("Bb3d", "0")

	// get url & set path
	u, err := url.ParseRequestURI(p.baseUrl)
	if err != nil {
		return err
	}
	u.Path = urlHttpCertServerSettings

	// make and do request
	req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// read body of response
	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// OK status?
	if resp.StatusCode != http.StatusOK {
		return errors.New("printer: failed to post set active cert form")
	}

	// find next CSRFToken
	csrfToken, err = parseBodyForCSRFToken(bodyBytes)
	if err != nil {
		return err
	}

	// submit confirmation (& reboot now)
	data = url.Values{}
	data.Set("pageid", "326")
	data.Set("CSRFToken", csrfToken)
	// 4 == do NOT activate other secure protos
	// 5 == DO activate other secure protos
	data.Set("http_page_mode", "5")

	// get url & set path
	u, err = url.ParseRequestURI(p.baseUrl)
	if err != nil {
		return err
	}
	u.Path = urlHttpCertServerSettings

	// make and do request
	req, err = http.NewRequest(http.MethodPost, u.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err = p.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// read body of response
	_, _ = io.Copy(io.Discard, resp.Body)

	// OK status?
	if resp.StatusCode != http.StatusOK {
		return errors.New("printer: failed to post set active cert form")
	}

	return nil
}
