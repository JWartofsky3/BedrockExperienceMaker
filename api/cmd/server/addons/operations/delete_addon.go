package operations

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/JWartofsky3/BedrockExperienceMaker/api/cmd/server/addons"
	"github.com/go-sql-driver/mysql"
)

func DeleteAddon(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := db.ExecContext(r.Context(), `DELETE FROM addons WHERE id = ?`, r.PathValue("name"))
		if err != nil {
			var mysqlError *mysql.MySQLError
			if errors.As(err, &mysqlError) && mysqlError.Number == 1451 {
				addons.WriteJSONError(w, "add-on is used by an experience pack", http.StatusBadRequest)
				return
			}
			addons.WriteJSONError(w, "could not delete add-on", http.StatusInternalServerError)
			return
		}
		count, err := result.RowsAffected()
		if err != nil {
			addons.WriteJSONError(w, "could not delete add-on", http.StatusInternalServerError)
			return
		}
		if count == 0 {
			addons.WriteJSONError(w, "add-on not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
