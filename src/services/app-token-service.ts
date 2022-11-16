import { createAppAuth } from "@octokit/auth-app";
import * as github from "@actions/github";

export class AppTokenService {
    async getToken() {
        const appId = Number(process.env.APP_ID);
        const privateKey = process.env.PRIVATE_KEY;
        const repository = process.env.GITHUB_REPOSITORY;

        if (!appId || !privateKey || !repository) {
            throw new Error("appId, privateKey, and repository are required");
        }

        return await this.generateToken(appId, privateKey, repository);
    }

    private async generateToken(appId: number, privateKey: string, repository: string) {
        const auth = createAppAuth({appId, privateKey});
        const {token: jwt} = await auth({type: "app"});
        const octokit = github.getOctokit(jwt);
        const [owner, repo] = repository.split("/");
        const {data: {id: installationId}} = await octokit.rest.apps.getRepoInstallation({owner, repo});
        const {token} = await auth({installationId, type: "installation"});

        return token;
    }
}
