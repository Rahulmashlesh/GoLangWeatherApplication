package webapp

/*func addToPoller(client client.HttpGetter, logger *slog.Logger, newZipcode string, unit string, metric metrics.Metrics, p *poller.Poller) {
	if acceptedZipCodes[newZipcode] {
		//already on the list, do nothing
	} else {
		acceptedZipCodes[newZipcode] = true
		p.Add(weather.NewCurrentWeather(client, logger, newZipcode, unit, metric))
	}
}

func zipCodeHandler(w http.ResponseWriter, r *http.Request, p *poller.Poller, client *client.OpenWeatherMapClient, logger *slog.Logger, unit string, metric metrics.Metrics) {

	fmt.Fprintln(w, "<html><body>")
	fmt.Fprintln(w, "<h1>Accepted Zip Codes</h1>")
	fmt.Fprintln(w, "<form method='post'>")
	fmt.Fprintln(w, "Zip Code: <input type='text' name='zipCode'>")
	fmt.Fprintln(w, "<input type='submit' value='Submit'>")
	fmt.Fprintln(w, "</form>")
	fmt.Fprintln(w, "</body></html>")

	if r.Method == http.MethodPost {
		r.ParseForm()
		newZipCode := r.FormValue("zipCode")
		if isValidZipCode(newZipCode) {
			addToPoller(client, logger, newZipCode, unit, metric, p)

		} else {
			fmt.Fprintln(w, "<p style='color:red;'>Invalid zip code:", newZipCode, "</p>")
		}
	}
	renderZipCodeList(w)
}

func isValidZipCode(zipCode string) bool {
	if zipCode == "0000" {
		return false
	}
	return true
}

func renderZipCodeList(w http.ResponseWriter) {
	fmt.Fprintln(w, "<ul>")
	for zipCode, _ := range acceptedZipCodes {
		fmt.Fprintf(w, "<li>%s</li>", zipCode)
	}
	fmt.Fprintln(w, "</ul>")
}

var (
	acceptedZipCodes = make(map[string]bool)
)

*/
