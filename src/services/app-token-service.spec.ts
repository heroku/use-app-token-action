import {
    AppAuthentication,
    InstallationAccessTokenAuthentication,
    InstallationAuthOptions,
    StrategyOptions
} from "@octokit/auth-app";
import {AppTokenService, AppTokenServiceInput} from "./app-token-service";

const input = {
    appId: "123456",
    privateKey: "fake-private-key",
    installationId: "99999999999999999999",
    repository: "fake-org/fake-github-repository",
} as AppTokenServiceInput;

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

beforeEach(() => {
    input.appId = "123456";
    input.privateKey = "fake-private-key";
    input.installationId = "99999999999999999999";
    input.repository = "fake-org/fake-github-repository";
});

afterEach(() => jest.restoreAllMocks());

describe("AppTokenService", () => {
    it("should generate a token successfully when an installationId is supplied", async () => {
        delete input.repository;

        const appTokenService = new AppTokenService(input);
        const expectedToken = "totally_fake_gh_installation_token";
        const actual = appTokenService.getToken()

        await expect(actual).resolves.not.toThrow();
        await expect(actual).resolves.toEqual(expectedToken);
    });

    it("should generate a token successfully when a repository is supplied", async () => {
        delete input.installationId

        const appTokenService = new AppTokenService(input);
        const expectedToken = "totally_fake_gh_installation_token";
        const actual = appTokenService.getToken()

        await expect(actual).resolves.not.toThrow();
        await expect(actual).resolves.toEqual(expectedToken);
    });

    it("should throw when there is no appId", async () => {
        input.appId = undefined as unknown as string;

        const appTokenService = () => new AppTokenService(input);

        await expect(appTokenService).toThrow(new Error("'appId' is required"));
    });

    it("should throw when there is no privateKey", async () => {
        input.privateKey = undefined as unknown as string;

        const appTokenService = () => new AppTokenService(input);

        await expect(appTokenService).toThrow(new Error("'privateKey' is required"));
    });

    it("should throw when there is no installationId and repository", async () => {
        delete input.installationId
        delete input.repository

        const appTokenService = () => new AppTokenService(input);

        await expect(appTokenService).toThrow(new Error("'installationId' or 'repository' is required"));
    });
});
