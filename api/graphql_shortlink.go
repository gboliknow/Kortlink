// graphql_shortlink.go
package api

import (
	"fmt"
	"kortlink/internal/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

func (s *ShortlinkService) graphqldeleteShortUrl() *graphql.Field {
    args := graphql.FieldConfigArgument{
        "shortURL": &graphql.ArgumentConfig{
            Type: graphql.NewNonNull(graphql.String),
        },
    }
    return &graphql.Field{
        Name:        "deleteShortUrl",
        Description: "Delete a short URL by passing the short URL",
        Type:        deleteShortUrlType,
        Args:        args,
        Resolve: func(p graphql.ResolveParams) (interface{}, error) {
            shortURL, ok := p.Args["shortURL"].(string)
            if !ok || shortURL == "" {
                return nil, fmt.Errorf("short URL is required")
            }

            // Check if the short URL exists
            _, err := s.store.GetOriginalURL(shortURL)
            if err != nil {
                return nil, fmt.Errorf("short URL not found")
            }

            // Delete the short URL
            err = s.store.DeleteShortURL(shortURL)
            if err != nil {
                return nil, fmt.Errorf("failed to delete short URL")
            }

            // Optionally delete from cache
            _ = s.cache.Delete(shortURL)

            return map[string]interface{}{
                "message":      "Short URL deleted successfully",
                "deletedCount": 1,
            }, nil
        },
    }
}

func (s *ShortlinkService) graphqlShortlinks(c *gin.Context) {
    rootQuery := graphql.NewObject(graphql.ObjectConfig{
        Name:   "Query",
        Fields: graphql.Fields{},
    })

    rootMutation := graphql.NewObject(graphql.ObjectConfig{
        Name: "Mutation",
        Fields: graphql.Fields{
            "deleteShortUrl": s.graphqldeleteShortUrl(),
        },
    })

    schema, _ := graphql.NewSchema(graphql.SchemaConfig{
        Query:    rootQuery,
        Mutation: rootMutation,
    })

    requestString := c.Query("q")
    if requestString == "" {
        var body map[string]interface{}
        if err := c.BindJSON(&body); err == nil {
            if query, ok := body["query"].(string); ok {
                requestString = query
            }
        }
    }

    res := graphql.Do(graphql.Params{
        Schema:        schema,
        RequestString: requestString,
    })

    if len(res.Errors) > 0 {
        utility.WriteJSON(c.Writer, http.StatusBadRequest, "GraphQL errors", res.Errors)
        return
    }

    utility.WriteJSON(c.Writer, http.StatusOK, "GraphQL fetched successfully", res)
}



var deleteShortUrlType = graphql.NewObject(graphql.ObjectConfig{
    Name: "deleteShortUrlType",
    Fields: graphql.Fields{
        "message": &graphql.Field{
            Type: graphql.String, // Success message
        },
        "deletedCount": &graphql.Field{
            Type: graphql.Int, // Count of deleted URLs (could be 1 for a single deletion)
        },
    },
})