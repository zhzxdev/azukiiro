package ranker

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/zhzxdev/azukiiro/client"
	"github.com/zhzxdev/azukiiro/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ParticipantResultSolution struct {
	Score       float64 `bson:"score"`
	Status      string  `bson:"status"`
	CompletedAt int     `bson:"completedAt"`
}

type ParticipantResult struct {
	SolutionCount  int                       `bson:"solutionCount"`
	LastSolutionId string                    `bson:"lastSolutionId"`
	LastSolution   ParticipantResultSolution `bson:"lastSolution"`
}

type Participant struct {
	Id        string                       `bson:"_id"`
	UserId    string                       `bson:"userId"`
	ContestId string                       `bson:"contestId"`
	Results   map[string]ParticipantResult `bson:"results"`
	UpdatedAt int                          `bson:"updatedAt"`
}

type ParticipantView struct {
	TotalScore int
	Raw        *Participant
}

type ByTotalScore []ParticipantView

func (a ByTotalScore) Len() int           { return len(a) }
func (a ByTotalScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTotalScore) Less(i, j int) bool { return a[i].TotalScore > a[j].TotalScore }

func Poll(ctx context.Context) (bool, error) {
	res, err := client.PollRanklist(ctx, &client.PollRanklistRequest{})
	if err != nil || res.TaskId == "" {
		return false, err
	}
	ctx = client.WithRanklistTask(ctx, res.TaskId, res.ContestId)
	// Sync solution list
	collection := db.Collection(ctx, fmt.Sprintf("contest-%v-participants", res.ContestId))
	now := time.Now().UnixMicro()
	since := 0
	participant := &Participant{}
	err = collection.FindOne(ctx, bson.D{}, options.FindOne().SetSort(bson.D{{Key: "updatedAt", Value: -1}})).Decode(participant)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return true, err
		}
	} else {
		since = participant.UpdatedAt
	}
	for {
		res, err := client.GetRanklistParticipants(ctx, since)
		if err != nil {
			return true, err
		}
		// TODO: optimize this logic
		if len(*res) == 0 {
			break
		}
		for _, value := range *res {
			_, err := collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: value.Id}}, bson.D{{Key: "$set", Value: value}}, options.Update().SetUpsert(true))
			if err != nil {
				return true, err
			}
		}
		last := (*res)[len(*res)-1]
		since = last.UpdatedAt
		if since >= int(now) {
			break
		}
	}
	var participants []ParticipantView
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return true, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var participant Participant
		err := cursor.Decode(&participant)
		if err != nil {
			return true, err
		}
		var totalScore int
		for _, value := range participant.Results {
			totalScore += int(value.LastSolution.Score)
		}
		participants = append(participants, ParticipantView{
			TotalScore: totalScore,
			Raw:        &participant,
		})
	}

	// Sort participants by total score
	sort.Sort(ByTotalScore(participants))

	// Sync ranklist
	problems, err := client.GetRanklistProblems(ctx)
	if err != nil {
		return true, err
	}

	var columns []*client.RanklistParticipantColumn
	columns = append(columns, &client.RanklistParticipantColumn{
		Name:        "Total",
		Description: "Total",
	})
	for _, value := range *problems {
		columns = append(columns, &client.RanklistParticipantColumn{
			Name:        value.Settings.Slug,
			Description: value.Title,
		})
	}

	ranklistMap := make(map[string]*client.Ranklist)
	for _, value := range res.Ranklists {
		// Currently there is no settings associated with ranklist, which is subject to change
		log.Println("Processing ranklist: ", value.Key)
		rank := 0
		var items []*client.RanklistParticipantItem
		for _, participant := range participants {
			var columns []*client.RanklistParticipantItemColumn
			var totalScore float64
			for _, problem := range *problems {
				if result, ok := participant.Raw.Results[problem.Id]; ok {
					totalScore += result.LastSolution.Score * problem.Settings.Score
				}
			}
			columns = append(columns, &client.RanklistParticipantItemColumn{
				Content: fmt.Sprintf("%v", totalScore),
			})
			for _, problem := range *problems {
				if result, ok := participant.Raw.Results[problem.Id]; ok {
					columns = append(columns, &client.RanklistParticipantItemColumn{
						Content: fmt.Sprintf("%v", result.LastSolution.Score),
					})
				} else {
					columns = append(columns, &client.RanklistParticipantItemColumn{
						Content: "0",
					})
				}
			}
			rank = rank + 1
			items = append(items, &client.RanklistParticipantItem{
				Rank:    rank,
				UserId:  participant.Raw.UserId,
				Columns: columns,
			})
		}

		ranklist := &client.Ranklist{
			Participant: &client.RanklistParticipant{
				Columns: columns,
				List:    items,
			},
			Metadata: &client.RanklistMetadata{
				GeneratedAt: int(now),
				Description: value.Name,
			},
		}
		ranklistMap[value.Key] = ranklist
	}
	if err = client.SaveRanklist(ctx, ranklistMap); err != nil {
		return true, err
	}
	return true, nil
}
