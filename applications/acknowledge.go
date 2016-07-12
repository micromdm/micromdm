package applications

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/micromdm/mdm"
	//	"github.com/micromdm/micromdm/device"
	"github.com/micromdm/micromdm/device"
)

func AppListPredicate(response mdm.Response) bool {
	fmt.Println("InstalledApplicationList Predicate")

	if response.RequestType == "InstalledApplicationList" {
		return true
	}

	if response.InstalledApplicationList != nil {
		return true
	}

	return false
}

func any(list []Application, predicate func(Application) bool) bool {
	for _, v := range list {
		if predicate(v) {
			return true
		}
	}

	return false
}

func AppListResponse(response mdm.Response, datastores map[string]interface{}) error {
	store, found := datastores["applications"]
	if !found {
		return errors.New("Do not have access to datastore for saving application information")
	}

	appsStore, ok := store.(Datastore)
	if !ok {
		return errors.New("could not acknowledge installed application list because the given datastore isnt an application datastore.")
	}

	dstore, found := datastores["devices"]
	if !found {
		return errors.New("Do not have access to datastore for retrieving device information")
	}

	deviceStore, ok := dstore.(device.Datastore)
	if !ok {
		return errors.New("could not acknowledge installed application list because the given device datastore isnt a device datastore.")
	}

	device, err := deviceStore.GetDeviceByUDID(response.UDID)
	if err != nil {
		return err
	}

	deviceApps, err := appsStore.GetApplicationsByDeviceUUID(device.UUID)
	if err != nil {
		return err
	}

	var uuids []string
	for _, app := range response.InstalledApplicationList {
		//var uuid string

		existingApps, err := appsStore.Applications(Name{app.Name}, Version{app.Version})
		if err != nil {
			return err
		}

		if len(existingApps) > 0 {
			existingApp := existingApps[0]
			uuid = existingApp.UUID
			uuids = append(uuids, existingApp.UUID)
		} else {
			dbApp := &Application{Name: app.Name}

			identifier := sql.NullString{}
			identifier.Scan(app.Identifier)
			dbApp.Identifier = identifier

			bundleSize := sql.NullInt64{}
			bundleSize.Scan(app.BundleSize)
			dbApp.BundleSize = bundleSize

			shortVersion := sql.NullString{}
			shortVersion.Scan(app.ShortVersion)
			dbApp.ShortVersion = shortVersion

			version := sql.NullString{}
			version.Scan(app.Version)
			dbApp.Version = version

			dynamicSize := sql.NullInt64{}
			dynamicSize.Scan(app.DynamicSize)
			dbApp.DynamicSize = dynamicSize

			isValidated := sql.NullBool{}
			isValidated.Scan(app.IsValidated)
			dbApp.IsValidated = isValidated

			uuid, err := appsStore.New(dbApp)
			if err != nil {
				return err
			}
			uuids = append(uuids, uuid)
		}

		if !any(deviceApps, app) {
			// App installed on device but not recorded

			//deviceApp := &DeviceApplication{
			//	DeviceUUID: device.UUID,
			//	ApplicationUUID: app.UUID,
			//}
		}
	}

	return nil
}
