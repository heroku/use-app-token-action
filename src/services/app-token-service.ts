import { createAppAuth } from "@octokit/auth-app";
import * as github from "@actions/github";


type NullableString = string | undefined;

export type AppTokenServiceInput = {
    appId: string,
    privateKey: string,
    installationId?: string,
    repository?: string
}

export class AppTokenService {
    constructor(private input: AppTokenServiceInput) {
        if (!input.appId) throw new Error("'appId' is required");
        if (!input.privateKey) throw new Error("'privateKey' is required");
        if (!input.installationId && !input.repository) throw new Error("'installationId' or 'repository' is required");
    }

    async getToken() {
        const auth = createAppAuth({appId: this.input.appId, privateKey: this.input.privateKey});
        const {token: jwt} = await auth({type: "app"});
        const installId = await this.getInstallationId(jwt)
        const {token} = await auth({installationId: installId, type: "installation"});

        return token;
    }

    private async getInstallationId(jwt: string) {
        if (this.input.installationId) return Number(this.input.installationId);

        const octokit = github.getOctokit(jwt);
        const [owner, repo] = this.input.repository?.split("/") || [];
        const {data: {id}} = await octokit.rest.apps.getRepoInstallation({owner, repo});

        return id;
    }
}
