"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.AppTokenService = void 0;
const auth_app_1 = require("@octokit/auth-app");
const github = __importStar(require("@actions/github"));
class AppTokenService {
    constructor(input) {
        this.input = input;
        if (!input.appId)
            throw new Error("'appId' is required");
        if (!input.privateKey)
            throw new Error("'privateKey' is required");
        if (!input.installationId && !input.repository)
            throw new Error("'installationId' or 'repository' is required");
    }
    async getToken() {
        const auth = (0, auth_app_1.createAppAuth)({ appId: this.input.appId, privateKey: this.input.privateKey });
        const { token: jwt } = await auth({ type: "app" });
        const installId = await this.getInstallationId(jwt);
        const { token } = await auth({ installationId: installId, type: "installation" });
        return token;
    }
    async getInstallationId(jwt) {
        var _a;
        if (this.input.installationId)
            return Number(this.input.installationId);
        const octokit = github.getOctokit(jwt);
        const [owner, repo] = ((_a = this.input.repository) === null || _a === void 0 ? void 0 : _a.split("/")) || [];
        const { data: { id } } = await octokit.rest.apps.getRepoInstallation({ owner, repo });
        return id;
    }
}
exports.AppTokenService = AppTokenService;
//# sourceMappingURL=app-token-service.js.map