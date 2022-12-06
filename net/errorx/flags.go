package errorx

import (
	"math/bits"
	"net/http"
	"sync"

	"github.com/jucardi/go-titan/errors"
)

var (
	flagsManager = &flagsHandler{}
)

func init() {
	// This sets a default map of errors.ErrFlag to an HTTP status code.
	//
	// See `Flags().SetMapping` for more information flags to HTTP status code mappings.
	_ = flagsManager.SetMapping(map[errors.ErrFlag]int{
		errors.FlagBadRequest:       http.StatusBadRequest,
		errors.FlagNotFound:         http.StatusNotFound,
		errors.FlagUnauthorized:     http.StatusUnauthorized,
		errors.FlagUnhandled:        http.StatusInternalServerError,
		errors.FlagNotImplemented:   http.StatusNotImplemented,
		errors.FlagOperationTimeout: http.StatusRequestTimeout,
		errors.FlagConflict:         http.StatusConflict,
	})
}

// IFlagsMappingHandler is the definition of an ErrFlags to HTTP Status code mapping manager
type IFlagsMappingHandler interface {
	// SetMapping allows to set a new mapping of ErrFlags (defined in the errors package of this
	// project) to HTTP status codes. This map helps the `Wrap` function determine what kind of HTTP
	// status code a newly wrapped error should have if the error has any flags set.
	//
	// The errors package of Titan has a predefined set of flags. However, a custom set of error flags
	// may be used instead. If using custom error flags, this map should be set to replace the default
	// with one that better represents HTTP status codes for the custom flags to be used.
	//
	// NOTE: Unique bit flags cannot have multiple HTTP status codes associated to them.
	//
	//   E.g: The following is a representation of a map[errors.ErrFlag]int, where the numbers on
	//        the left are the binary representation of a flag value and the numbers on the right
	//        are the HTTP status codes
	//
	//        {
	//          1000: 404
	//           100: 401
	//             1: 500
	//            11: 403
	//        }
	//
	//        In that example, that mapping attempts to assign two HTTP status codes to the bit 0x1.
	//        This will result in an error when the mapping is processed because mapping an HTTP status
	//        code to the bit 0x1 would be ambiguous (500 and 403)
	//
	// NOTE: For cases where multiple flag bits can point to the same HTTP status code, when doing a
	//       reverse conversion, the status code with the greater value will be the one used for this
	//       conversion
	SetMapping(m map[errors.ErrFlag]int) error

	// ToFlags attempts to convert an HTTP status code to errors.ErrFlag by using the defined mapping
	// of this instance
	ToFlags(httpStatus int) errors.ErrFlag

	// ToStatus attempts to convert errors.ErrFlag to an HTTP status code by using the defined mapping
	// of this instance
	ToStatus(flags errors.ErrFlag) int

	// IsStatusInFlags indicates whether a set of flags are mappable to the specified HTTP status code
	// This function is useful for cases where an `errors.ErrFlags` value has multiple flags which would
	// make the exact flags match fail in the flags to status map.
	IsStatusInFlags(flags errors.ErrFlag, status int) bool
}

// Flags returns the errors.ErrFlag to HTTP status codes mapping manager.
func Flags() IFlagsMappingHandler {
	return flagsManager
}

type flagsHandler struct {
	flagToStatus map[errors.ErrFlag]int
	statusToFlag map[int]errors.ErrFlag
	mux          *sync.RWMutex
}

func (f *flagsHandler) SetMapping(m map[errors.ErrFlag]int) error {
	if len(m) == 0 {
		return errors.New("map is empty, no changes made")
	}

	newFlagToStatus := map[errors.ErrFlag]int{}
	newStatusToFlag := map[int]errors.ErrFlag{}

	for flag, status := range m {
		if fCount := bits.OnesCount(uint(flag)); fCount == 0 {
			// Ignoring empty flags
			continue
		} else if fCount == 1 {
			if err := f.setMapping(newFlagToStatus, newStatusToFlag, flag, status); err != nil {
				return err
			}
		} else {
			fls := f.splitFlags(flag)
			for _, fl := range fls {
				if err := f.setMapping(newFlagToStatus, newStatusToFlag, fl, status); err != nil {
					return err
				}
			}
		}
	}

	f.mux.Lock()
	defer f.mux.Unlock()

	f.flagToStatus = newFlagToStatus
	f.statusToFlag = newStatusToFlag

	return nil
}

func (f *flagsHandler) ToFlags(httpStatus int) errors.ErrFlag {
	if ret, ok := f.statusToFlag[httpStatus]; ok {
		return ret
	}
	return 0
}

func (f *flagsHandler) ToStatus(flags errors.ErrFlag) int {
	if ret, ok := f.flagToStatus[flags]; ok {
		return ret
	}
	return 0
}

func (f *flagsHandler) IsStatusInFlags(flags errors.ErrFlag, status int) bool {
	fl, ok := f.statusToFlag[status]
	if !ok {
		return false
	}
	return fl&flags == flags
}

func (f *flagsHandler) splitFlags(flags errors.ErrFlag) []errors.ErrFlag {
	var (
		shiftingBit errors.ErrFlag
		ret         []errors.ErrFlag
	)

	for i := int8(0); i < 32; i++ {
		shiftingBit = 1 << i
		if flags&shiftingBit == shiftingBit {
			ret = append(ret, shiftingBit)
		}
	}
	return ret
}

func (f *flagsHandler) setMapping(newFlagToStatus map[errors.ErrFlag]int, newStatusToFlag map[int]errors.ErrFlag, flag errors.ErrFlag, status int) error {
	if existingStatus, exists := newFlagToStatus[flag]; exists && existingStatus != status {
		return errors.Format("unable to set mapping, flag with value %d has ambiguous HTTP status codes (%d, %d)", flag, status, existingStatus)
	}

	newFlagToStatus[flag] = status

	if existingFlag, ok := newStatusToFlag[status]; ok {
		newStatusToFlag[status] = existingFlag | flag
	} else {
		newStatusToFlag[status] = flag

	}
	return nil
}
