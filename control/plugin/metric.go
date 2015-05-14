package plugin

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/intelsdi-x/pulse/core/cdata"
	"github.com/intelsdi-x/pulse/pkg/logger"
)

const (
	// These are our built-in content types for plugins

	// PulseAll the wildcard for accepting all pulse content types
	PulseAllContentType = "pulse.*"
	// PulseGOB pulse metrics serialized into go binary format
	PulseGOBContentType = "pulse.gob"
	// PulseJSON pulse metrics serialized into json
	PulseJSONContentType = "pulse.json"
	// PulseProtoBuff pulse metrics serialized into protocol buffers
	// PulseProtoBuff = "pulse.pb" // TO BE IMPLEMENTED
)

// Represents a metric type. Only used within plugins and across plugin calls.
// Converted to core.MetricType before being used within modules.
type PluginMetricType struct {
	Namespace_          []string              `json:"namespace"`
	LastAdvertisedTime_ time.Time             `json:"last_advertised_time"`
	Version_            int                   `json:"version"`
	Config_             *cdata.ConfigDataNode `json:"config"`
	Data_               interface{}           `json:"data"`
}

// // PluginMetricType Constructor
func NewPluginMetricType(namespace []string, data interface{}) *PluginMetricType {
	return &PluginMetricType{
		Namespace_: namespace,
		Data_:      data,
	}
}

// Returns the namespace.
func (p PluginMetricType) Namespace() []string {
	return p.Namespace_
}

// Returns the last time this metric type was received from the plugin.
func (p PluginMetricType) LastAdvertisedTime() time.Time {
	return p.LastAdvertisedTime_
}

// Returns the namespace.
func (p PluginMetricType) Version() int {
	return p.Version_
}

// Config returns the map of config data for this metric
func (p PluginMetricType) Config() *cdata.ConfigDataNode {
	return p.Config_
}

func (p PluginMetricType) Data() interface{} {
	return p.Data_
}

func (p *PluginMetricType) AddData(data interface{}) {
	p.Data_ = data
}

// MarshallMetricTypes returns a []byte containing a serialized version of []PluginMetricType using the content type provided.
func MarshallPluginMetricTypes(contentType string, metrics []PluginMetricType) ([]byte, string, error) {
	// If we have an empty slice we return an error
	if len(metrics) == 0 {
		es := fmt.Sprintf("attempt to marshall empty slice of metrics: %s", contentType)
		logger.Error("marshal-metric-types", es)
		return nil, "", errors.New(es)
	}
	// Switch on content type
	switch contentType {
	case PulseAllContentType, PulseGOBContentType:
		// NOTE: A Pulse All wildcard will result in GOB
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(metrics)
		if err != nil {
			logger.Error("marshal-metric-types", err.Error())
			return nil, "", err
		}
		// contentType := PulseGOBContentType
		return buf.Bytes(), PulseGOBContentType, nil
	case PulseJSONContentType:
		// Serialize into JSON
		b, err := json.Marshal(metrics)
		if err != nil {
			logger.Error("marshal-metric-types", err.Error())
			return nil, "", err
		}
		return b, PulseJSONContentType, nil
	default:
		// We don't recognize this content type. Log and return error.
		es := fmt.Sprintf("invlaid pulse content type: %s", contentType)
		logger.Error("marshal-metric-types", es)
		return nil, "", errors.New(es)
	}
}

// UnmarshallPluginMetricTypes takes a content type and []byte payload and returns a []PluginMetricType
func UnmarshallPluginMetricTypes(contentType string, payload []byte) ([]PluginMetricType, error) {
	switch contentType {
	case PulseGOBContentType:
		var metrics []PluginMetricType
		r := bytes.NewBuffer(payload)
		err := gob.NewDecoder(r).Decode(&metrics)
		if err != nil {
			logger.Error("unmarshal-metric-types", err.Error())
			return nil, err
		}
		return metrics, nil
	case PulseJSONContentType:
		var metrics []PluginMetricType
		err := json.Unmarshal(payload, &metrics)
		if err != nil {
			logger.Error("unmarshal-metric-types", err.Error())
			return nil, err
		}
		return metrics, nil
	default:
		// We don't recognize this content type as one we can unmarshal. Log and return error.
		es := fmt.Sprintf("invlaid pulse content type for unmarshalling: %s", contentType)
		logger.Error("unmarshal-metric-types", es)
		return nil, errors.New(es)
	}
}

// SwapPluginMetricContentType swaps a payload with one content type to another one.
func SwapPluginMetricContentType(contentType, requestedContentType string, payload []byte) ([]byte, string, error) {
	metrics, err1 := UnmarshallPluginMetricTypes(contentType, payload)
	if err1 != nil {
		logger.Error("swap-content-type", err1.Error())
		return nil, "", err1
	}
	newPayload, newContentType, err2 := MarshallPluginMetricTypes(requestedContentType, metrics)
	if err2 != nil {
		logger.Error("swap-content-type", err2.Error())
		return nil, "", err2
	}
	return newPayload, newContentType, nil
}
