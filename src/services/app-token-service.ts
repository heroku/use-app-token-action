import { createAppAuth } from "@octokit/auth-app";
import * as github from "@actions/github";


type NullableString = string | undefined;

export class AppTokenService {
    async getToken() {
        const appId = Number(process.env.APP_ID);
        const privateKey = process.env.PRIVATE_KEY;
        const installationId = process.env.INSTALLATION_ID;
        const repository = process.env.GITHUB_REPOSITORY;

        if (!appId) throw new Error("APP_ID is required");
        if (!privateKey) throw new Error("PRIVATE_KEY is required");
        if (!installationId && !repository) throw new Error("INSTALLATION_ID or GITHUB_REPOSITORY is required");

        return await this.generateToken(appId, privateKey, installationId, repository);
    }

    private async generateToken(appId: number, privateKey: string, installationId: NullableString, repository: NullableString) {
        const auth = createAppAuth({appId, privateKey});
        const {token: jwt} = await auth({type: "app"});
        const installId = await this.getInstallationId(jwt, installationId, repository)
        console.log(`installId: ${installId}`)
        const {token} = await auth({installationId: installId, type: "installation"});

        return token;
    }

    private async getInstallationId(jwt: string, installationId: NullableString, repository: NullableString) {
        if (installationId) return Number(installationId);

        const octokit = github.getOctokit(jwt);
        const [owner, repo] = repository?.split("/") || [];
        const {data: {id}} = await octokit.rest.apps.getRepoInstallation({owner, repo});

        return id;
    }
}
