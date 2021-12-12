import { Children } from "react"
import { Route } from "react-location"

export const routes: Route[] = [
    {
        path: '/',
    },
    {
        path: "login"
    },
    {
        path: "/apps",
        children: [
            {
                path: "/"
            },
            {
                path: "/:appId",
                loader: async ({ params: { teamId } }) => ({
                    team: await fetch(`/api/teams/${teamId}`),
                  }),
            }
        ]
    }
]
