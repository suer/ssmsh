package instance

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

const instanceIDPrefix = "i-"

type Candidate struct {
	InstanceID string
	Name       string
	State      string
}

type Selector func(name string, candidates []Candidate) (string, error)

type NotFoundError struct {
	Name string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("no running instance found for Name tag %q", e.Name)
}

func Resolve(ctx context.Context, client EC2API, name string, selector Selector) (string, error) {
	if strings.HasPrefix(name, instanceIDPrefix) {
		return name, nil
	}

	candidates, err := findByNameTag(ctx, client, name)
	if err != nil {
		return "", err
	}

	switch len(candidates) {
	case 0:
		return "", &NotFoundError{Name: name}
	case 1:
		return candidates[0].InstanceID, nil
	default:
		return selector(name, candidates)
	}
}

func findByNameTag(ctx context.Context, client EC2API, name string) ([]Candidate, error) {
	input := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{Name: aws.String("tag:Name"), Values: []string{name}},
			{Name: aws.String("instance-state-name"), Values: []string{"running"}},
		},
	}

	var candidates []Candidate
	paginator := ec2.NewDescribeInstancesPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("describe instances: %w", err)
		}
		for _, r := range page.Reservations {
			for _, inst := range r.Instances {
				candidates = append(candidates, toCandidate(inst))
			}
		}
	}
	return candidates, nil
}

func toCandidate(inst types.Instance) Candidate {
	c := Candidate{InstanceID: aws.ToString(inst.InstanceId)}
	if inst.State != nil {
		c.State = string(inst.State.Name)
	}
	for _, tag := range inst.Tags {
		if aws.ToString(tag.Key) == "Name" {
			c.Name = aws.ToString(tag.Value)
		}
	}
	return c
}
