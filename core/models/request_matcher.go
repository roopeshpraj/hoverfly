package models

import (
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/util"
)

type RequestFieldMatchers struct {
	ExactMatch    *string
	XmlMatch      *string
	XpathMatch    *string
	JsonMatch     *string
	JsonPathMatch *string
	RegexMatch    *string
	GlobMatch     *string
}

func NewRequestFieldMatchersFromView(matchers *v2.RequestFieldMatchersView) *RequestFieldMatchers {
	if matchers == nil {
		return nil
	}

	return &RequestFieldMatchers{
		ExactMatch:    matchers.ExactMatch,
		XmlMatch:      matchers.XmlMatch,
		XpathMatch:    matchers.XpathMatch,
		JsonMatch:     matchers.JsonMatch,
		JsonPathMatch: matchers.JsonPathMatch,
		RegexMatch:    matchers.RegexMatch,
		GlobMatch:     matchers.GlobMatch,
	}
}

func (this RequestFieldMatchers) BuildView() *v2.RequestFieldMatchersView {
	return &v2.RequestFieldMatchersView{
		ExactMatch:    this.ExactMatch,
		XmlMatch:      this.XmlMatch,
		XpathMatch:    this.XpathMatch,
		JsonMatch:     this.JsonMatch,
		JsonPathMatch: this.JsonPathMatch,
		RegexMatch:    this.RegexMatch,
		GlobMatch:     this.GlobMatch,
	}
}

type RequestMatcherResponsePair struct {
	RequestMatcher RequestMatcher
	Response       ResponseDetails
}

func NewRequestMatcherResponsePairFromView(view *v2.RequestResponsePairViewV2) *RequestMatcherResponsePair {
	if view.Request.Query != nil && view.Request.Query.ExactMatch != nil {
		sortedQuery := util.SortQueryString(*view.Request.Query.ExactMatch)
		view.Request.Query.ExactMatch = &sortedQuery
	}

	return &RequestMatcherResponsePair{
		RequestMatcher: RequestMatcher{
			Path:        NewRequestFieldMatchersFromView(view.Request.Path),
			Method:      NewRequestFieldMatchersFromView(view.Request.Method),
			Destination: NewRequestFieldMatchersFromView(view.Request.Destination),
			Scheme:      NewRequestFieldMatchersFromView(view.Request.Scheme),
			Query:       NewRequestFieldMatchersFromView(view.Request.Query),
			Body:        NewRequestFieldMatchersFromView(view.Request.Body),
			Headers:     view.Request.Headers,
		},
		Response: NewResponseDetailsFromResponse(view.Response),
	}
}

func (this *RequestMatcherResponsePair) BuildView() v2.RequestResponsePairViewV2 {

	var path, method, destination, scheme, query, body *v2.RequestFieldMatchersView

	if this.RequestMatcher.Path != nil {
		path = this.RequestMatcher.Path.BuildView()
	}

	if this.RequestMatcher.Method != nil {
		method = this.RequestMatcher.Method.BuildView()
	}

	if this.RequestMatcher.Destination != nil {
		destination = this.RequestMatcher.Destination.BuildView()
	}

	if this.RequestMatcher.Scheme != nil {
		scheme = this.RequestMatcher.Scheme.BuildView()
	}

	if this.RequestMatcher.Query != nil {
		query = this.RequestMatcher.Query.BuildView()
	}

	if this.RequestMatcher.Body != nil {
		body = this.RequestMatcher.Body.BuildView()
	}

	return v2.RequestResponsePairViewV2{
		Request: v2.RequestDetailsViewV2{
			Path:        path,
			Method:      method,
			Destination: destination,
			Scheme:      scheme,
			Query:       query,
			Body:        body,
			Headers:     this.RequestMatcher.Headers,
		},
		Response: this.Response.ConvertToResponseDetailsView(),
	}
}

type RequestMatcher struct {
	Path        *RequestFieldMatchers
	Method      *RequestFieldMatchers
	Destination *RequestFieldMatchers
	Scheme      *RequestFieldMatchers
	Query       *RequestFieldMatchers
	Body        *RequestFieldMatchers
	Headers     map[string][]string
}

func (this RequestMatcher) BuildRequestDetailsFromExactMatches() *RequestDetails {
	if this.Body == nil || this.Body.ExactMatch == nil ||
		this.Destination == nil || this.Destination.ExactMatch == nil ||
		this.Method == nil || this.Method.ExactMatch == nil ||
		this.Path == nil || this.Path.ExactMatch == nil ||
		this.Query == nil || this.Query.ExactMatch == nil ||
		this.Scheme == nil || this.Scheme.ExactMatch == nil {
		return nil
	}

	return &RequestDetails{
		Body:        *this.Body.ExactMatch,
		Destination: *this.Destination.ExactMatch,
		Headers:     this.Headers,
		Method:      *this.Method.ExactMatch,
		Path:        *this.Path.ExactMatch,
		Query:       *this.Query.ExactMatch,
		Scheme:      *this.Scheme.ExactMatch,
	}

}
