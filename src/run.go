package multiqc

import (
	"os"
	"os/exec"
	"path"
	"time"

	codeclarity "github.com/CodeClarityCE/utility-types/codeclarity_db"
	exceptionManager "github.com/CodeClarityCE/utility-types/exceptions"
	"github.com/uptrace/bun"

	"github.com/parithera/plugin-multiqc/src/types"
	"github.com/parithera/plugin-multiqc/src/utils/output_generator"
)

// Start is the main entry point for executing the MultiQC script.
// It takes the source code directory and a CodeClarity DB connection as input.
func Start(sourceCodeDir string, codeclarityDB *bun.DB) types.Output {
	return ExecuteScript(sourceCodeDir) // Calls the ExecuteScript function to perform the analysis.
}

// ExecuteScript executes the MultiQC script in the specified directory.
// It handles the creation of the output directory and the execution of the MultiQC command.
func ExecuteScript(sourceCodeDir string) types.Output {
	start := time.Now() // Records the start time of the analysis.

	outputPath := path.Join(sourceCodeDir, "multiqc") // Defines the output directory for MultiQC.
	os.MkdirAll(outputPath, os.ModePerm)              // Creates the output directory if it doesn't exist.

	fastqcPath := path.Join(sourceCodeDir, "fastqc") // Defines the path to the FastQC directory.
	fastpPath := path.Join(sourceCodeDir, "fastp")   // Defines the path to the FastP directory.
	starPath := path.Join(sourceCodeDir, "STAR")     // Defines the path to the STAR directory.

	args := []string{"-o", outputPath, fastqcPath, fastpPath, starPath} // Constructs the command-line arguments for MultiQC.
	// Run Rscript in sourceCodeDir
	cmd := exec.Command("multiqc", args...) // Creates an exec.Command object to run the MultiQC command.
	_, err := cmd.CombinedOutput()          // Executes the command and captures the combined output (stdout and stderr).
	if err != nil {
		codeclarity_error := exceptionManager.Error{
			Private: exceptionManager.ErrorContent{
				Description: err.Error(),
				Type:        exceptionManager.GENERIC_ERROR,
			},
			Public: exceptionManager.ErrorContent{
				Description: "The script failed to execute",
				Type:        exceptionManager.GENERIC_ERROR,
			},
		}
		return generate_output(start, nil, codeclarity.FAILURE, []exceptionManager.Error{codeclarity_error}) // Returns an output with the error information.
	}

	return generate_output(start, "done", codeclarity.SUCCESS, []exceptionManager.Error{}) // Returns an output indicating successful completion.
}

// generate_output formats the analysis results into a types.Output struct.
// It includes the analysis timing, status, and any errors that occurred.
func generate_output(start time.Time, data any, status codeclarity.AnalysisStatus, errors []exceptionManager.Error) types.Output {
	formattedStart, formattedEnd, delta := output_generator.GetAnalysisTiming(start) // Gets the analysis timing information.

	output := types.Output{
		Result: types.Result{
			Data: data, // The analysis data.
		},
		AnalysisInfo: types.AnalysisInfo{
			Errors: errors, // Any errors that occurred during the analysis.
			Time: types.Time{
				AnalysisStartTime: formattedStart, // The start time of the analysis.
				AnalysisEndTime:   formattedEnd,   // The end time of the analysis.
				AnalysisDeltaTime: delta,          // The duration of the analysis.
			},
			Status: status, // The status of the analysis (success or failure).
		},
	}
	return output
}
