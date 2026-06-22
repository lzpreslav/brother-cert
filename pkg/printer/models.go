package printer

import (
	"fmt"
	"sort"
	"strings"
)

// Brother's Web Based Management uses opaque, per-model form field names (e.g.
// "Bb23") for the certificate import / activate / delete pages. The names are
// stable for a given model+firmware family but differ between models, so each
// supported model carries its own field map here. Everything else in this
// package (login, CSRF token, cert-list parsing) is discovered dynamically and
// is model-independent.

// importFields names the fields of the certificate import form
// (/net/security/certificate/import.html).
type importFields struct {
	pageID string // hidden "pageid" value identifying the form
	empty1 string // first opaque hidden field, submitted empty
	empty2 string // second opaque hidden field, submitted empty
	file   string // multipart file field carrying the p12
	passwd string // p12 password field (we use a password-less p12, so empty)
}

// activateFields names the fields of the HTTP server settings form
// (net/net/certificate/http.html) used to select the active certificate.
type activateFields struct {
	pageID   string // hidden "pageid" value
	dropdown string // <select> naming the active certificate id
	// protocols re-asserts every protocol toggle the firmware exposes on this
	// page. An absent checkbox is treated as "off", so to change only the
	// certificate we must resubmit each protocol with its enabled value. The
	// map is name->value (value "1" to enable, "" / "0" where the firmware
	// expects an empty/zero companion field).
	protocols map[string]string
	// httpPageMode: 4 == keep other secure protocols as-is, 5 == also activate
	// them. Submitted in the confirmation step.
	httpPageMode string
}

// deleteFields names the fields of the certificate delete form's first step
// (/net/security/certificate/delete.html). The confirmation step's hidden
// fields are scraped dynamically (see parseHiddenInputs) and need no mapping.
type deleteFields struct {
	pageID string // hidden "pageid" value
	empty1 string // first opaque hidden field, submitted empty
	empty2 string // second opaque hidden field, submitted empty
}

// model bundles the field maps for one printer model.
type model struct {
	name     string
	imp      importFields
	activate activateFields
	del      deleteFields
}

// models is the registry of supported printer models, keyed by the value
// passed to --model. Add a new model by capturing its import/activate/delete
// form fields (e.g. from a browser HAR) and registering it here.
var models = map[string]model{
	// The model the upstream gregtwallace/brother-cert project targets.
	"MFC-L2710DW": {
		name: "MFC-L2710DW",
		imp: importFields{
			pageID: "390",
			empty1: "B8ea",
			empty2: "B8f8",
			file:   "B820",
			passwd: "B821",
		},
		activate: activateFields{
			pageID:   "326",
			dropdown: "B903",
			protocols: map[string]string{
				"B86c": "1",
				"B87e": "1",
			},
			httpPageMode: "5",
		},
		del: deleteFields{
			pageID: "383",
			empty1: "B8ea",
			empty2: "B8fc",
		},
	},

	// Brother MFC-L2750DW laser MFP.
	"MFC-L2750DW": {
		name: "MFC-L2750DW",
		imp: importFields{
			pageID: "395",
			empty1: "Bb0a",
			empty2: "Bb18",
			file:   "Ba40",
			passwd: "Ba41",
		},
		activate: activateFields{
			pageID:   "326",
			dropdown: "Bb23",
			protocols: map[string]string{
				"Ba8c":         "1", // Web Based Management HTTPS (443)
				"Ba8d":         "1", // Web Based Management HTTP  (80)
				"Ba9e":         "1", // IPP HTTPS (443)
				"ipp_ssl_used": "",  // IPP secure helper (submitted empty)
				"Ba9f":         "1", // IPP HTTP (80)
				"Baa0":         "1", // IPP HTTP (631)
				"Ba7d":         "1", // Web Services HTTP
				"Bb20":         "",
				"Bb21":         "",
				"Bb3d":         "0",
			},
			httpPageMode: "5",
		},
		del: deleteFields{
			pageID: "388",
			empty1: "Bb0a",
			empty2: "Bb1c",
		},
	},
}

// lookupModel returns the field map for the named model, or an error listing
// the supported models if the name is unknown.
func lookupModel(name string) (model, error) {
	m, ok := models[name]
	if !ok {
		return model{}, fmt.Errorf("printer: unsupported model %q (supported: %s)", name, strings.Join(supportedModels(), ", "))
	}
	return m, nil
}

// supportedModels returns the registered model names, sorted.
func supportedModels() []string {
	names := make([]string, 0, len(models))
	for n := range models {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
