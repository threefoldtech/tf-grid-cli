package router

import "fmt"

func generateProjectName(projectName string) string {
	return fmt.Sprintf("%s.tfgrid.cli", projectName)
}
