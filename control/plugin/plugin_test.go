/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2015 Intel Coporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package plugin

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPluginType(t *testing.T) {
	Convey(".String()", t, func() {
		Convey("it returns the correct string representation", func() {
			p := PluginType(0)
			So(p.String(), ShouldEqual, "collector")
		})
	})
}

func TestMetricType(t *testing.T) {
	Convey("MetricType", t, func() {
		now := time.Now()
		m := &PluginMetricType{
			Namespace_:          []string{"foo", "bar"},
			LastAdvertisedTime_: now,
		}
		Convey("New", func() {
			So(m, ShouldHaveSameTypeAs, &PluginMetricType{})
		})
		Convey("Get Namespace", func() {
			So(m.Namespace(), ShouldResemble, []string{"foo", "bar"})
		})
		Convey("Get LastAdvertisedTimestamp", func() {
			So(m.LastAdvertisedTime().Unix(), ShouldBeGreaterThan, now.Unix()-2)
			So(m.LastAdvertisedTime().Unix(), ShouldBeLessThan, now.Unix()+2)
		})
	})
}

func TestArg(t *testing.T) {
	Convey("NewArg", t, func() {
		arg := NewArg("/tmp/pulse/plugin.log")
		So(arg, ShouldNotBeNil)
	})
}

func TestPlugin(t *testing.T) {
	a := []string{PulseAllContentType}
	b := []string{PulseGOBContentType}
	Convey("Start", t, func() {
		mockPluginMeta := NewPluginMeta("test", 1, CollectorPluginType, a, b)
		var mockPluginArgs string = "{\"PluginLogPath\": \"/var/tmp/pulse_plugin.log\"}"
		err, rc := Start(mockPluginMeta, new(MockPlugin), mockPluginArgs)
		So(err, ShouldBeNil)
		So(rc, ShouldEqual, 0)
	})
	Convey("Start with invalid args", t, func() {
		mockPluginMeta := NewPluginMeta("test", 1, CollectorPluginType, a, b)
		var mockPluginArgs string = ""
		err, rc := Start(mockPluginMeta, new(MockPlugin), mockPluginArgs)
		So(err, ShouldNotBeNil)
		So(rc, ShouldNotEqual, 0)
	})
	Convey("Bad accept type", t, func() {
		a := []string{"wat"}
		So(func() {
			NewPluginMeta("test", 1, CollectorPluginType, a, b)
		}, ShouldPanicWith, "Bad accept content type [test] for [1] [wat]")
	})
	Convey("Bad return type", t, func() {
		b := []string{"wat"}
		So(func() {
			NewPluginMeta("test", 1, CollectorPluginType, a, b)
		}, ShouldPanicWith, "Bad return content type [test] for [1] [wat]")
	})
}
