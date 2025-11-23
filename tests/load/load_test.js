import http from "k6/http";
import { check, sleep } from "k6";

export let options = {
    vus: 10,
    duration: "5m",
};

const BASE = "http://localhost:8080";

const TEAMS_COUNT = 30;
const USERS_PER_TEAM = 15;

function uuidv4() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
        const r = Math.random() * 16 | 0;
        const v = c === 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}

export function setup() {
    const initialTeams = [];
    for (let t = 1; t <= TEAMS_COUNT; t++) {
        const teamName = `team-${uuidv4()}`;
        const members = [];
        for (let i = 1; i <= USERS_PER_TEAM; i++) {
            members.push({
                user_id: uuidv4(),
                username: `User${uuidv4()}`,
                is_active: true,
            });
        }

        const teamResp = http.post(`${BASE}/team/add`, JSON.stringify({
            team_name: teamName,
            members: members
        }), { headers: { "Content-Type": "application/json" }});

        check(teamResp, {
            "setup team created or exists": (r) => r.status === 201 || r.status === 400
        });

        initialTeams.push({ teamName, members });
    }
    return initialTeams;
}

export default function (initialTeams) {
    const vu = __VU;
    const iter = __ITER;

    const team = initialTeams[Math.floor(Math.random() * initialTeams.length)];
    const members = team.members;

    const author = members[0];

    const isActive = Math.random() > 0.5;
    http.post(`${BASE}/users/setIsActive`, JSON.stringify({
        user_id: author.user_id,
        is_active: isActive
    }), { headers: { "Content-Type": "application/json" }});

    const prId = uuidv4();
    const prResp = http.post(`${BASE}/pullRequest/create`, JSON.stringify({
        pull_request_id: prId,
        pull_request_name: `Test PR ${prId}`,
        author_id: author.user_id
    }), { headers: { "Content-Type": "application/json" }});

    check(prResp, {
        "PR created or conflict": (r) => r.status === 201 || r.status === 400 || r.status === 409
    });

    if (Math.random() > 0.3) {
        const mergeResp = http.post(`${BASE}/pullRequest/merge`, JSON.stringify({
            pull_request_id: prId
        }), { headers: { "Content-Type": "application/json" }});

        check(mergeResp, {
            "PR merged or already merged": (r) => r.status === 200 || r.status === 400 || r.status === 404
        });
    }

    const statsResp = http.get(`${BASE}/stats/reviewers`);
    check(statsResp, {
        "stats fetched": (r) => r.status === 200 || r.status === 204
    });

    sleep(Math.random() * 0.5 + 0.5);
}
