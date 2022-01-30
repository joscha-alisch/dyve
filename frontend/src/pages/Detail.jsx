import {useParams} from "react-router-dom";
import {useEffect, useState} from "react";
import axios from "axios";
import Page from "../components/base/pages/page/page";
import API from "../api";
import DetailPage from "../components/base/pages/detailPage/detailpage";
import Box from "../components/base/box/box";
import {prettyJ} from "../helpers/pretty";
import RoutingCard from "../components/cards/routing/routingcard";
import InstancesCard from "../components/cards/instances/instancescard";

export const AppDetail = () => <DetailPage
    parentRoute={"/apps"}
    parent={"Apps"}
    detailApi={API.Apps.Get}
    useLive={API.Apps.Live}
    render={(detail, live) => <div>
        <RoutingCard routing={live.routing}/>
        <InstancesCard instances={live.instances} />
    </div>}
/>

export const TeamDetail = () => <DetailPage
    parentRoute={"/teams"}
    parent={"Teams"}
    detailApi={API.Teams.Get}
    deleteEnabled={true}
    editRoute={(id) => "/teams/" + id + "/edit"}
    render={(team) => <>
        <Box>{team.description}</Box>
    </>}
/>