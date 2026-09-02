package instance

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type fakeEC2 struct {
	output *ec2.DescribeInstancesOutput
	err    error
}

func (f *fakeEC2) DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	return f.output, f.err
}

func failSelector(t *testing.T) Selector {
	return func(name string, candidates []Candidate) (string, error) {
		t.Fatal("selector should not be called")
		return "", nil
	}
}

func TestResolve_InstanceIDPrefix(t *testing.T) {
	client := &fakeEC2{err: errors.New("should not be called")}
	id, err := Resolve(context.Background(), client, "i-05239d460d79ad389", failSelector(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "i-05239d460d79ad389" {
		t.Fatalf("got %q, want %q", id, "i-05239d460d79ad389")
	}
}

func TestResolve_NotFound(t *testing.T) {
	client := &fakeEC2{output: &ec2.DescribeInstancesOutput{}}
	_, err := Resolve(context.Background(), client, "test-ec2-terraform-minimum", failSelector(t))
	var notFound *NotFoundError
	if !errors.As(err, &notFound) {
		t.Fatalf("got %v, want NotFoundError", err)
	}
}

func TestResolve_SingleMatch(t *testing.T) {
	client := &fakeEC2{output: &ec2.DescribeInstancesOutput{
		Reservations: []types.Reservation{
			{Instances: []types.Instance{
				{InstanceId: aws.String("i-05239d460d79ad389"), Tags: []types.Tag{
					{Key: aws.String("Name"), Value: aws.String("test-ec2-terraform-minimum")},
				}},
			}},
		},
	}}
	id, err := Resolve(context.Background(), client, "test-ec2-terraform-minimum", failSelector(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "i-05239d460d79ad389" {
		t.Fatalf("got %q, want %q", id, "i-05239d460d79ad389")
	}
}

func TestResolve_MultipleMatches_UsesSelector(t *testing.T) {
	client := &fakeEC2{output: &ec2.DescribeInstancesOutput{
		Reservations: []types.Reservation{
			{Instances: []types.Instance{
				{InstanceId: aws.String("i-aaa"), Tags: []types.Tag{{Key: aws.String("Name"), Value: aws.String("dup")}}},
				{InstanceId: aws.String("i-bbb"), Tags: []types.Tag{{Key: aws.String("Name"), Value: aws.String("dup")}}},
			}},
		},
	}}
	selectorCalled := false
	selector := func(name string, candidates []Candidate) (string, error) {
		selectorCalled = true
		if name != "dup" {
			t.Fatalf("got name %q, want %q", name, "dup")
		}
		if len(candidates) != 2 {
			t.Fatalf("got %d candidates, want 2", len(candidates))
		}
		return "i-bbb", nil
	}
	id, err := Resolve(context.Background(), client, "dup", selector)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !selectorCalled {
		t.Fatal("selector was not called")
	}
	if id != "i-bbb" {
		t.Fatalf("got %q, want %q", id, "i-bbb")
	}
}

func TestResolve_DescribeInstancesError(t *testing.T) {
	client := &fakeEC2{err: errors.New("boom")}
	_, err := Resolve(context.Background(), client, "test-ec2-terraform-minimum", failSelector(t))
	if err == nil {
		t.Fatal("expected error")
	}
}
