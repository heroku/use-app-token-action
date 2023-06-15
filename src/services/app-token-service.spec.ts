import {
    AppAuthentication,
    InstallationAccessTokenAuthentication,
    InstallationAuthOptions,
    StrategyOptions
} from "@octokit/auth-app";
import {AppTokenService} from "./app-token-service";

const env = process.env;

const octokitRestMocks = {
    apps: {
        getRepoInstallation: (params: { owner: string, repo: string }) => Promise.resolve({
            status: 200,
            data: {
                id: 99999999999999999999,
            }
        })
    }
};

jest.mock("@actions/core", () => ({
    __esModule: true,
    ...jest.requireActual("@actions/core")
}));

jest.mock("@actions/github", () => ({
    __esModule: true,
    ...jest.requireActual("@actions/github"),
    getOctokit: () => ({rest: octokitRestMocks})
}));

beforeEach(() => {
    process.env = {
        ...process.env,
        APP_ID: "123456",
        PRIVATE_KEY: "fake-private-key",
        INSTALLATION_ID: "99999999999999999999",
        GITHUB_REPOSITORY: "fake-org/fake-github-repository",
    }
})

afterEach(() => {
    jest.restoreAllMocks();
    process.env = env
});

jest.mock("@octokit/auth-app", () => ({
    __esModule: true,
    ...jest.requireActual("@octokit/auth-app"),
    createAppAuth: (params: StrategyOptions) => {
        return (params: AppAuthentication | InstallationAuthOptions) => {
            let resolvedValue: AppAuthentication | InstallationAccessTokenAuthentication = null as unknown as AppAuthentication;
            const now = new Date().toISOString();

            if (params.type === "app") {
                resolvedValue = {
                    appId: params.appId,
                    type: "app",
                    token: "totally_fake_gh_jwt",
                    expiresAt: now
                } as AppAuthentication;
            } else if (params.type === "installation") {
                resolvedValue = {
                    installationId: Number(params.installationId),
                    type: "token",
                    tokenType: "installation",
                    token: "totally_fake_gh_installation_token",
                    expiresAt: now
                } as InstallationAccessTokenAuthentication
            }

            return Promise.resolve(resolvedValue)
        };
    }
}));

describe("AppTokenService", () => {
    const appTokenService = new AppTokenService()

    it("should generate a token successfully when an installationId is supplied", async () => {
        delete process.env.GITHUB_REPOSITORY;

        const expectedToken = "totally_fake_gh_installation_token";
        const actual = appTokenService.getToken()

        await expect(actual).resolves.not.toThrow();
        await expect(actual).resolves.toEqual(expectedToken);
    });

    it("should generate a token successfully when a repository is supplied", async () => {
        delete process.env.INSTALLATION_ID

        const expectedToken = "totally_fake_gh_installation_token";
        const actual = appTokenService.getToken()

        await expect(actual).resolves.not.toThrow();
        await expect(actual).resolves.toEqual(expectedToken);
    });

    it("should throw when there is no appId", async () => {
        delete process.env.APP_ID;

        const actual = appTokenService.getToken()

        await expect(actual).rejects.toThrow(new Error("APP_ID is required"));
    });

    it("should throw when there is no privateKey", async () => {
        delete process.env.PRIVATE_KEY

        const actual = appTokenService.getToken()

        await expect(actual).rejects.toThrow(new Error("PRIVATE_KEY is required"));
    });

    it("should throw when there is no installationId and repository", async () => {
        delete process.env.INSTALLATION_ID
        delete process.env.GITHUB_REPOSITORY

        const actual = appTokenService.getToken()

        await expect(actual).rejects.toThrow(new Error("INSTALLATION_ID or GITHUB_REPOSITORY is required"));
    });
});
