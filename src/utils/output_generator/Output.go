package output_generator

import (
	"time"

	sbomTypes "github.com/CodeClarityCE/plugin-sbom-javascript/src/types/sbom/js"
	codeclarity "github.com/CodeClarityCE/utility-types/codeclarity_db"
	exceptionManager "github.com/CodeClarityCE/utility-types/exceptions"
)

// GetAnalysisTiming calculates the start time, end time, and elapsed time of the analysis.
// It takes the analysis start time as a parameter and returns the formatted start time, end time, and elapsed time in seconds.
func GetAnalysisTiming(start time.Time) (string, string, float64) {
	end := time.Now()
	elapsed := time.Since(start)
	return start.Local().String(), end.Local().String(), elapsed.Seconds()
}

// WriteFailureOutput writes the failure output for the analysis.
// It sets the analysis status to 'FAILURE', updates the analysis timing information (start time, end time, and duration),
// and retrieves/sets any errors that occurred during the analysis.
// The updated output object is then returned.
func WriteFailureOutput(output sbomTypes.Output, start time.Time) sbomTypes.Output {
	output.AnalysisInfo.Status = codeclarity.FAILURE
	formattedStart, formattedEnd, delta := GetAnalysisTiming(start)
	output.AnalysisInfo.Time.AnalysisStartTime = formattedStart
	output.AnalysisInfo.Time.AnalysisEndTime = formattedEnd
	output.AnalysisInfo.Time.AnalysisDeltaTime = delta

	output.AnalysisInfo.Errors = exceptionManager.GetErrors()

	return output
}
