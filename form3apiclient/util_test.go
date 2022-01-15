package form3apiclient

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const someValidRelativePath = "some/path"
const someValidAbsolutePath = "http://example.com"

var _ = Describe("util", func() {
	Context("join() function", func() {
		DescribeTable("joins absoluth and relative urls",
			func(baseUrl, path, expectedResult string) {
				actualResultUrl, err := join(baseUrl, path)

				Expect(err).NotTo(HaveOccurred())
				Expect(actualResultUrl).To(Equal(expectedResult))
			},
			EntryDescription(`"%s" join "%s" is "%s"`),
			Entry(nil, "http://example.com", "some/path", "http://example.com/some/path"),
			Entry(nil, "http://example.com", "/some/path", "http://example.com/some/path"),
			Entry(nil, "http://example.com/", "some/path", "http://example.com/some/path"),
			Entry(nil, "http://example.com/", "/some/path", "http://example.com/some/path"),
			Entry(nil, "http://example.com/", "/some/path/", "http://example.com/some/path"),

			Entry(nil, "http://example.com/v1", "some/path", "http://example.com/v1/some/path"),
			Entry(nil, "http://example.com/v1", "/some/path", "http://example.com/v1/some/path"),
			Entry(nil, "http://example.com/v1/", "/some/path", "http://example.com/v1/some/path"),
			Entry(nil, "http://example.com/v1/", "some/path", "http://example.com/v1/some/path"),
			Entry(nil, "http://example.com/v1/", "some/path/", "http://example.com/v1/some/path"),
		)

		DescribeTable("invalid parameters cause error",
			func(baseUrl, path, expectedError string) {
				actualResultUrl, err := join(baseUrl, path)

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(fmt.Errorf(expectedError)))
				Expect(actualResultUrl).To(BeZero())
			},
			EntryDescription(`"%s" join "%s" causes error "%s"`),
			Entry(nil, "example.com/v1", someValidRelativePath, "baseAbsoluteUrl must be absolute"),
			Entry(nil, "http://example.com/v1?param=value", someValidRelativePath, "baseAbsoluteUrl with query is not supported"),
			Entry(nil, "http://example.com/v1#fragment", someValidRelativePath, "baseAbsoluteUrl with fragment is not supported"),

			Entry(nil, someValidAbsolutePath, "/some/path?param=value", "relativeUrl with query is not supported"),
			Entry(nil, someValidAbsolutePath, "/some/path#fragment", "relativeUrl with fragment is not supported"),
		)
	})
})
