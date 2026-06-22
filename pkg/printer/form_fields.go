package printer

import "regexp"

var (
	hiddenInputRegex = regexp.MustCompile(`(?s)<input[^>]*type="hidden"[^>]*>`)
	inputNameRegex   = regexp.MustCompile(`name="([^"]*)"`)
	inputValueRegex  = regexp.MustCompile(`(?s)value="(.*?)"`)
)

// parseHiddenInputs returns a name->value map of every hidden <input> field in
// the supplied HTML. It is used to forward a confirmation form's hidden fields
// (whose randomized empty-field names vary by page) without hardcoding them.
func parseHiddenInputs(body []byte) map[string]string {
	out := map[string]string{}
	for _, tag := range hiddenInputRegex.FindAll(body, -1) {
		nameMatch := inputNameRegex.FindSubmatch(tag)
		if nameMatch == nil {
			continue
		}
		value := ""
		if v := inputValueRegex.FindSubmatch(tag); v != nil {
			value = string(v[1])
		}
		out[string(nameMatch[1])] = value
	}
	return out
}
