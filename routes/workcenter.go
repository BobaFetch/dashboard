package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bobafetch/dashboard/data"
)

type Workcenter struct {
	WCNREF  string `json:"wcRef"`
	WCNDESC string `json:"wcDesc"`
}

type Workorder struct {
	RUNREF      string  `json:"runRef"`
	RUNRTNUM    string  `json:"runRtNum"`
	RUNNO       string  `json:"runNo"`
	RUNQTY      string  `json:"runQty"`
	SOCUST      string  `json:"soCust"`
	COMMENTS    string  `json:"comments"`
	RUNPRIORITY string  `json:"runPriority"`
	OPSCHEDDATE *string `json:"opSchedDate"`
	DATEDIFF    string  `json:"opDateDiff"`
	WCNDESC     string  `json:"wcDesc"`
}

func RegisterWorkcenterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/workcenter", getWorkcenters)
	mux.HandleFunc("GET /api/workcenter/{ref}", getWorkcenter)
}

func getWorkcenters(w http.ResponseWriter, r *http.Request) {
	db := data.GetDB()

	query := "SELECT WCNREF, WCNDESC FROM WcntTable"
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var workcenters []Workcenter
	for rows.Next() {
		var wc Workcenter
		err := rows.Scan(&wc.WCNREF, &wc.WCNDESC)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Row scan error", http.StatusInternalServerError)
			return
		}
		workcenters = append(workcenters, wc)
	}

	jsonData, err := json.Marshal(workcenters)
	if err != nil {
		http.Error(w, "JSON conversion error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func getWorkcenter(w http.ResponseWriter, r *http.Request) {
	ref := r.PathValue("ref")
	fmt.Printf("Fetching workcenter: %s\n", ref)
	db := data.GetDB()

	query := `
	       SELECT DISTINCT
		       RTRIM(RUNREF) RUNREF,
		       RTRIM(RUNRTNUM) RUNRTNUM,
		       RUNNO,
		       RUNQTY,
		       RTRIM(SOCUST) SOCUST,
		       ISNULL(AGPMCOMMENTS, '') COMMENTS,
		       RUNPRIORITY,
		       OPSCHEDDATE,
		       ISNULL((SELECT DATEDIFF(MINUTE, (SELECT TOP 1 OPCOMPDATE FROM RnopTable WHERE OPREF=RUNREF AND OPRUN=RUNNO AND OPCOMPLETE IS NOT NULL ORDER BY OPCOMPDATE DESC), GETDATE())), 0) DATEDIFF,
		       RTRIM(WCNDESC) WCNDESC
	       FROM RunsTable
		       INNER JOIN RnopTable ON OPREF=RUNREF AND OPRUN=RUNNO AND RUNOPCUR=OPNO
		       INNER JOIN PartTable ON PARTREF=RUNREF
		       INNER JOIN RnalTable ON RUNREF=RAREF AND RUNNO=RARUN
		       INNER JOIN SohdTable ON SONUMBER=RASO
		       INNER JOIN WcntTable ON OPCENTER=WCNREF
		       LEFT OUTER JOIN AgcmTable ON AGPART=RUNRTNUM AND AGRUN=RUNNO
		       LEFT OUTER JOIN SoitTable ON ITPART=PARTREF AND ITSO=RASO
		       WHERE OPCENTER = @ref AND OPCOMPLETE=0
		       ORDER BY RUNPRIORITY, OPSCHEDDATE`

	rows, err := db.Query(query, sql.Named("ref", ref))
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	if rows != nil {
		defer rows.Close()
	}

	var workorders []Workorder
	for rows != nil && rows.Next() {
		var wo Workorder
		err := rows.Scan(
			&wo.RUNREF,
			&wo.RUNRTNUM,
			&wo.RUNNO,
			&wo.RUNQTY,
			&wo.SOCUST,
			&wo.COMMENTS,
			&wo.RUNPRIORITY,
			&wo.OPSCHEDDATE,
			&wo.DATEDIFF,
			&wo.WCNDESC,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Row scan error", http.StatusInternalServerError)
			return
		}
		workorders = append(workorders, wo)
	}

	jsonData, err := json.Marshal(workorders)
	if err != nil {
		http.Error(w, "JSON conversion error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
