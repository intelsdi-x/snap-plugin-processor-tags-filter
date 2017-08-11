// +build small

/*
http://www.apache.org/licenses/LICENSE-2.0.txt

Copyright 2017 Intel Corporation

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

package tagsfilter

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

func TestTagsFilter(t *testing.T) {
	tfp := NewTFProcessor()

	Convey("Create tag processor", t, func() {
		Convey("So tfp should not be nil", func() {
			So(tfp, ShouldNotBeNil)
		})
		Convey("So tfp should be of type TFProcessor", func() {
			So(tfp, ShouldHaveSameTypeAs, &TFProcessor{})
		})
		Convey("tfp.GetConfigPolicy should return a config policy", func() {
			configPolicy, _ := tfp.GetConfigPolicy()
			Convey("So config policy should be a plugin.ConfigPolicy", func() {
				So(configPolicy, ShouldHaveSameTypeAs, plugin.ConfigPolicy{})
			})
		})
	})

	Convey("Test TagsFilter Processor", t, func() {
		Convey("Process metrics with one allowed value per tag", func() {
			config := plugin.Config{
				"test1.allow": "fooval",
				"test2.allow": "barval",
			}
			metrics := []plugin.Metric{
				{
					Namespace: plugin.NewNamespace("foo"),
					Data:      456,
					Tags:      map[string]string{"test1": "fooval", "test2": "barval"},
				},
				{
					Namespace: plugin.NewNamespace("foo"),
					Data:      789,
					Tags:      map[string]string{"test1": "value", "test2": "barval"},
				},
				{
					Namespace: plugin.NewNamespace("foo"),
					Data:      123,
					Tags:      map[string]string{"test1": "barval", "test2": "fooval"},
				},
			}
			mts, err := tfp.Process(metrics, config)
			So(mts, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(mts, ShouldHaveLength, 2)
			So(mts[0].Data, ShouldEqual, 456)
			So(mts[1].Data, ShouldEqual, 789)
		})

		Convey("Process metrics with more than one allowed value per tag", func() {
			config := plugin.Config{
				"test1.allow": "fooval,barval",
				"test2.allow": "bazval",
			}
			metrics := []plugin.Metric{
				{
					Namespace: plugin.NewNamespace("foo"),
					Data:      456,
					Tags:      map[string]string{"test1": "fooval"},
				},
				{
					Namespace: plugin.NewNamespace("foo"),
					Data:      789,
					Tags:      map[string]string{"test1": "barval", "test2": "something"},
				},
				{
					Namespace: plugin.NewNamespace("foo"),
					Data:      258,
					Tags:      map[string]string{"test1": "bezval", "test2": "bazval"},
				},
				{
					Namespace: plugin.NewNamespace("foo"),
					Data:      123,
					Tags:      map[string]string{"test1": "some other value", "test3": "another value"},
				},
			}
			mts, err := tfp.Process(metrics, config)
			So(mts, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(mts, ShouldHaveLength, 3)
			So(mts[0].Data, ShouldEqual, 456)
			So(mts[1].Data, ShouldEqual, 789)
			So(mts[2].Data, ShouldEqual, 258)
		})

		Convey("Process metrics with no metrics matching allowing rules", func() {
			config := plugin.Config{
				"baz.allow": "bazval",
			}
			metrics := []plugin.Metric{
				{
					Namespace: plugin.NewNamespace("foo"),
					Data:      123,
					Tags:      map[string]string{"test1": "fooval"},
				},
				{
					Namespace: plugin.NewNamespace("foo"),
					Data:      456,
					Tags:      map[string]string{"test1": "barval", "test2": "something"},
				},
			}
			mts, err := tfp.Process(metrics, config)
			So(mts, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(mts, ShouldHaveLength, 0)
		})

		Convey("Process metrics with one denied value per tag", func() {
			config := plugin.Config{
				"test1.allow": "fooval",
				"test2.allow": "bazval",
				"test2.deny":  "barval",
			}
			metrics := []plugin.Metric{
				{
					Namespace: plugin.NewNamespace("foo"),
					Data:      456,
					Tags:      map[string]string{"test1": "fooval", "test2": "barval"},
				},
				{
					Namespace: plugin.NewNamespace("foo"),
					Data:      789,
					Tags:      map[string]string{"test1": "barval", "test2": "fooval"},
				},
				{
					Namespace: plugin.NewNamespace("foo"),
					Data:      123,
					Tags:      map[string]string{"test1": "otherval", "test2": "bazval"},
				},
			}
			mts, err := tfp.Process(metrics, config)
			So(mts, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(mts, ShouldHaveLength, 1)
			So(mts[0].Data, ShouldEqual, 123)
		})

		Convey("Process metrics with no metrics matching denying rules", func() {
			config := plugin.Config{
				"test1.allow": "fooval",
				"test1.deny":  "otherval",
				"test2.allow": "bazval",
				"test2.deny":  "barval",
			}
			metrics := []plugin.Metric{
				{
					Namespace: plugin.NewNamespace("foo"),
					Data:      456,
					Tags:      map[string]string{"test1": "fooval", "test2": "barval"},
				},
				{
					Namespace: plugin.NewNamespace("foo"),
					Data:      789,
					Tags:      map[string]string{"test1": "barval", "test2": "fooval"},
				},
				{
					Namespace: plugin.NewNamespace("foo"),
					Data:      123,
					Tags:      map[string]string{"test1": "otherval", "test2": "bazval"},
				},
			}
			mts, err := tfp.Process(metrics, config)
			So(mts, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(mts, ShouldHaveLength, 0)
		})

		Convey("Process metrics with no tag name in config key", func() {
			config := plugin.Config{
				"baz.allow": "test",
				".allow":    "test2",
			}
			metrics := []plugin.Metric{
				{
					Namespace: plugin.NewNamespace("foo"),
					Data:      123,
					Tags:      map[string]string{"test1": "fooval"},
				},
			}
			mts, err := tfp.Process(metrics, config)
			So(mts, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})

		Convey("Process metrics with no allowed value for rule", func() {
			config := plugin.Config{
				"baz.allow": "",
				"foo.allow": "test",
			}
			metrics := []plugin.Metric{
				{
					Namespace: plugin.NewNamespace("foo"),
					Data:      123,
					Tags:      map[string]string{"test1": "fooval"},
				},
			}
			mts, err := tfp.Process(metrics, config)
			So(mts, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})

		Convey("Process metrics with one of allowed values empty", func() {
			config := plugin.Config{
				"baz.allow": "bazval,",
				"foo.allow": "test",
			}
			metrics := []plugin.Metric{
				{
					Namespace: plugin.NewNamespace("foo"),
					Data:      123,
					Tags:      map[string]string{"baz": "bazval"},
				},
			}
			mts, err := tfp.Process(metrics, config)
			So(mts, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})

		Convey("Process metrics with no rules", func() {
			config := plugin.Config{
				"somekey": "somevalue",
			}
			metrics := []plugin.Metric{
				{
					Namespace: plugin.NewNamespace("foo"),
					Data:      123,
					Tags:      map[string]string{"baz": "bazval"},
				},
				{
					Namespace: plugin.NewNamespace("bar"),
					Data:      456,
					Tags:      map[string]string{"baz": "bazval"},
				},
			}
			mts, err := tfp.Process(metrics, config)
			So(mts, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(mts, ShouldHaveLength, 2)
			So(mts[0].Data, ShouldEqual, 123)
			So(mts[1].Data, ShouldEqual, 456)
		})

		Convey("Process metrics with empty config", func() {
			config := plugin.Config{}
			metrics := []plugin.Metric{
				{
					Namespace: plugin.NewNamespace("foo"),
					Data:      123,
					Tags:      map[string]string{"baz": "bazval"},
				},
			}
			mts, err := tfp.Process(metrics, config)
			So(mts, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(mts, ShouldHaveLength, 1)
			So(mts[0].Data, ShouldEqual, 123)
		})

		Convey("Process metrics with wrong config", func() {
			config := plugin.Config{
				"foo.allow": 12,
			}
			metrics := []plugin.Metric{
				{
					Namespace: plugin.NewNamespace("foo"),
					Data:      123,
					Tags:      map[string]string{"baz": "bazval"},
				},
			}
			mts, err := tfp.Process(metrics, config)
			So(mts, ShouldBeNil)
			So(err.Error(), ShouldEqual, "config item is not a string")
			So(mts, ShouldHaveLength, 0)
		})
	})
}
