package api

import "time"

// Application encapsulates the info of an application.
type Application struct {
	Index     string    // string
	Mode      string    // mode
	Version   string    // version
	CreatedAt time.Time // created time
	UpdatedAt time.Time // updated time
}

// GetApplication returns the application info for the given app identified by index.
//
func GetApplication(index string) (*Application, error) {
	return nil, nil
}

// GetVersion returns the version info for the given app identified by index.
//
func GetVersion(index string) string {
	app, _ := GetApplication(index)

	if app != nil {
		return app.Version
	}
	return ""
}
