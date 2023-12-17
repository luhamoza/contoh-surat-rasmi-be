package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jung-kurt/gofpdf"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/pdf", handlePDFGeneration).Methods("POST", "OPTIONS")

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})

	fmt.Println("Starting server on port 8080")
	http.ListenAndServe(":8080", handlers.CORS(originsOk, headersOk, methodsOk)(router))
}

func handlePDFGeneration(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Could not parse form data", http.StatusInternalServerError)
		return
	}

	// Get form values
	yourName := r.FormValue("yourName")
	yourAddress := r.FormValue("yourAddress")
	date := r.FormValue("date")
	employerName := r.FormValue("employerName")
	companyName := r.FormValue("companyName")
	companyAddress := r.FormValue("companyAddress")
	yourPosition := r.FormValue("yourPosition")
	endWorkDate := r.FormValue("endWorkDate")
	yourSignature := r.FormValue("yourSignature")

	// Initialize a new pdf
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)

	// Add formal letter
	formalLetter := fmt.Sprintf(
		"%s\n%s\n%s\n\n%s\n%s\n%s\n\nTuan/Puan,\n\nPerletakan Jawatan\n\nSaya ingin memaklumkan bahawa saya memutuskan untuk meletakkan jawatan dari %s di %s, berkuatkuasa pada %s.\n\nSaya bersyukur atas peluang yang diberikan selama saya berkhidmat di sini. Selama tempoh notis, saya berjanji untuk menyelesaikan semua tugas dan tanggungjawab dengan sebaiknya dan akan membantu dalam proses transisi kepada pekerja yang akan menggantikan saya.\n\nTerima kasih atas sokongan dan bimbingan yang diberikan sepanjang saya berkhidmat di %s. Saya berharap dapat terus menjaga hubungan profesional yang baik di masa hadapan.\n\nHormat saya,\n\n%s\n%s",
		yourName, yourAddress, date, employerName, companyName, companyAddress, yourPosition, companyName, endWorkDate, companyName, yourSignature, yourName)

	// Set left, top margins. 10mm each side (A4: 210mm wide).
	pdf.SetLeftMargin(10)
	pdf.SetTopMargin(10)

	// MultiCell(width, height, content, border [, align [, fill [, link]]])
	// width = 0 (full width), height = 50 (a rough approximation for this content).
	pdf.MultiCell(190, 9, formalLetter, "", "", false)

	// Set headers
	w.Header().Set("Content-Disposition", "attachment; filename=FormalLetter.pdf")
	w.Header().Set("Content-Type", "application/pdf")

	// Output PDF
	err = pdf.Output(w)
	if err != nil {
		http.Error(w, "Could not generate PDF", http.StatusInternalServerError)
	}
}
