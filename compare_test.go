package deep_test

import (
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hottestseason/deep"
)

func TestDiff(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want any
		got  any
		diff string
	}{
		{
			name: "1:same",
			want: map[string]any{
				"foo": "bar",
			},
			got: map[string]any{
				"foo": "bar",
			},
		},
		{
			name: "1:different",
			want: map[string]any{
				"foo": "bar",
			},
			got: map[string]any{},
			diff: `{map[string]any}["foo"]:
	-: bar (string)
	+: <invalid reflect.Value>
`,
		},
		{
			name: "2:same",
			want: map[string]any{
				"foo": []string{"bar", "baz"},
				"qux": deep.SuperHashOf(map[string]any{
					"quux": "corge",
				}),
			},
			got: map[string]any{
				"foo": []any{"bar", "baz"},
				"qux": map[string]any{
					"quux":  "corge",
					"graut": "garply",
				},
			},
		},
		{
			name: "2:different",
			want: map[string]any{
				"foo": []string{"bar"},
				"qux": deep.SuperHashOf(map[string]any{
					"corge": "quux",
				}),
			},
			got: map[string]any{
				"foo": []any{"bar", "baz"},
				"qux": map[string]any{
					"quux":  "corge",
					"graut": "garply",
				},
			},
			diff: `{map[string]any}["foo"]:
	-: [bar] ([]string len:1)
	+: [bar baz] ([]interface {} len:2)

{map[string]any}["qux"]:
	-: is super hash of map[corge:quux] (deep.superHashOfMatcher)
	+: map[graut:garply quux:corge] (map[string]interface {})
`,
		},
		{
			name: "3:same",
			want: deep.SuperHashOf(map[string]any{
				"baz": []any{
					deep.SuperHashOf(map[string]any{
						"qux": deep.InAnyOrder([]string{"a", "b", "c"}),
					}),
				},
			}),
			got: map[string]any{
				"foo": "bar",
				"baz": []map[string]any{
					{
						"qux":  []string{"c", "a", "b"},
						"quux": "corge",
					},
				},
			},
		},
		{
			name: "3:different",
			want: deep.SuperHashOf(map[string]any{
				"baz": []any{
					deep.SuperHashOf(map[string]any{
						"qux": deep.InAnyOrder([]string{"a", "b", "c"}),
					}),
				},
			}),
			got: map[string]any{
				"foo": "bar",
				"baz": []map[string]any{
					{
						"qux":  []string{"c", "a", "d"},
						"quux": "corge",
					},
				},
			},
			diff: `{map[string]any}.(map[string]any)["baz"]{[]map[string]any}{map[string]any}.(map[string]any)["qux"]{[]string}:
	-: has the same elements as [a b c] (deep.inAnyOrderMatcher[string])
	+: [c a d] ([]string len:3)
`,
		},
		{
			name: "4:same",
			want: deep.JSON(deep.SuperHashOf(map[string]any{
				"baz": deep.InAnyOrder([]any{
					map[string]any{
						"garply": gomock.Any(),
					},
					deep.SuperHashOf(map[string]any{
						"qux": 2,
					}),
				}),
			})),
			got: `{ "foo": "bar", "baz": [{ "qux": 2, "corge": "graut" }, { "garply": "waldo" }] }`,
		},
		{
			name: "4:different",
			want: deep.JSON(deep.SuperHashOf(map[string]any{
				"baz": deep.InAnyOrder([]any{
					map[string]any{
						"garply": gomock.Any(),
						"fred":   "plugh",
					},
					deep.SuperHashOf(map[string]any{
						"qux": 2,
					}),
				}),
			})),
			got: `{ "foo": "bar", "baz": [{ "qux": 2, "corge": "graut" }, { "garply": "waldo" }] }`,
			diff: `{string}{map[string]any}.(map[string]any)["baz"]{[]any}:
	-: has the same elements as [map[fred:plugh garply:is anything] is super hash of map[qux:2]] (deep.inAnyOrderMatcher[interface {}])
	+: [map[corge:graut qux:2] map[garply:waldo]] ([]interface {} len:3)
`,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if diff := strings.ReplaceAll(deep.Diff(tt.want, tt.got), "\u00A0", " "); diff != tt.diff {
				t.Errorf("deep.Diff(%v, %v) = \n%v\nwant\n%v", tt.want, tt.got, diff, tt.diff)
			}
		})
	}
}
