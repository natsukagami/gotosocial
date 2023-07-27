// Code generated by astool. DO NOT EDIT.

package streams

import (
	typeemoji "github.com/superseriousbusiness/activity/streams/impl/toot/type_emoji"
	typehashtag "github.com/superseriousbusiness/activity/streams/impl/toot/type_hashtag"
	typeidentityproof "github.com/superseriousbusiness/activity/streams/impl/toot/type_identityproof"
	vocab "github.com/superseriousbusiness/activity/streams/vocab"
)

// NewTootEmoji creates a new TootEmoji
func NewTootEmoji() vocab.TootEmoji {
	return typeemoji.NewTootEmoji()
}

// NewTootHashtag creates a new TootHashtag
func NewTootHashtag() vocab.TootHashtag {
	return typehashtag.NewTootHashtag()
}

// NewTootIdentityProof creates a new TootIdentityProof
func NewTootIdentityProof() vocab.TootIdentityProof {
	return typeidentityproof.NewTootIdentityProof()
}
