package form3apiclient

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	someValidRelativePath = "some/path"
	someValidAbsolutePath = "http://example.com"
)

func noDescEntry(args ...interface{}) TableEntry {
	return Entry(nil, args...)
}

var _ = Describe("util", func() {
	Context("join() function", func() {
		DescribeTable("joins absoluth and relative urls",
			func(baseURL, path, expectedResult string) {
				actualResultURL, err := join(baseURL, path)

				Expect(err).NotTo(HaveOccurred())
				Expect(actualResultURL).To(Equal(expectedResult))
			},
			EntryDescription(`"%s" join "%s" is "%s"`),
			noDescEntry("http://example.com", "some/path", "http://example.com/some/path"),
			noDescEntry("http://example.com", "/some/path", "http://example.com/some/path"),
			noDescEntry("http://example.com/", "some/path", "http://example.com/some/path"),
			noDescEntry("http://example.com/", "/some/path", "http://example.com/some/path"),
			noDescEntry("http://example.com/", "/some/path/", "http://example.com/some/path"),

			noDescEntry("http://example.com/v1", "some/path", "http://example.com/v1/some/path"),
			noDescEntry("http://example.com/v1", "/some/path", "http://example.com/v1/some/path"),
			noDescEntry("http://example.com/v1/", "/some/path", "http://example.com/v1/some/path"),
			noDescEntry("http://example.com/v1/", "some/path", "http://example.com/v1/some/path"),
			noDescEntry("http://example.com/v1/", "some/path/", "http://example.com/v1/some/path"),
		)

		DescribeTable("invalid parameters cause error",
			func(baseURL, path, expectedError string) {
				actualResultURL, err := join(baseURL, path)

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(fmt.Errorf(expectedError)))
				Expect(actualResultURL).To(BeZero())
			},
			EntryDescription(`"%s" join "%s" causes error "%s"`),
			noDescEntry("example.com/v1", someValidRelativePath, "baseAbsoluteURL must be absolute"),
			noDescEntry(
				"http://example.com/v1?param=value",
				someValidRelativePath,
				"baseAbsoluteURL with query is not supported"),
			noDescEntry(
				"http://example.com/v1#fragment",
				someValidRelativePath,
				"baseAbsoluteURL with fragment is not supported"),

			noDescEntry(someValidAbsolutePath, "/some/path?param=value", "relativeURL with query is not supported"),
			noDescEntry(someValidAbsolutePath, "/some/path#fragment", "relativeURL with fragment is not supported"),
		)
	})
})
