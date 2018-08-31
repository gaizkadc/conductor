//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the server instance

package server

import (
    "fmt"
    "net/http"
    "time"

    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/provider/configstorage"
    "github.com/daishogroup/system-model/server/config"
    log "github.com/sirupsen/logrus"
    "github.com/gorilla/mux"

    "github.com/daishogroup/system-model/provider/appdescstorage"
    "github.com/daishogroup/system-model/provider/appinststorage"
    "github.com/daishogroup/system-model/provider/clusterstorage"
    "github.com/daishogroup/system-model/provider/networkstorage"
    "github.com/daishogroup/system-model/provider/userstorage"
    "github.com/daishogroup/system-model/server/app"
    "github.com/daishogroup/system-model/server/cluster"
    "github.com/daishogroup/system-model/server/network"
    "github.com/daishogroup/system-model/provider/nodestorage"
    "github.com/daishogroup/system-model/server/node"
    "github.com/daishogroup/system-model/server/dump"
    "github.com/daishogroup/system-model/server/backup"
    "github.com/daishogroup/system-model/server/info"
    "github.com/daishogroup/system-model/server/user"
    "github.com/daishogroup/system-model/server/session"
    "github.com/daishogroup/system-model/provider/accessstorage"
    "github.com/daishogroup/system-model/server/access"
    "github.com/daishogroup/system-model/provider/passwordstorage"
    "github.com/daishogroup/system-model/server/password"
    "github.com/daishogroup/system-model/provider/oauthstorage"
    "github.com/daishogroup/system-model/provider/credentialsstorage"
    "github.com/daishogroup/system-model/provider/sessionstorage"
    "github.com/daishogroup/system-model/server/oauth"
    "github.com/daishogroup/system-model/errors"
    "github.com/daishogroup/system-model/server/credentials"
    "github.com/daishogroup/dhttp"
)

// Logger for the Service struct.
var loggerService = log.WithField(
    "package", "api",
).WithField(
    "struct", "Service",
).WithField(
    "file", "service.go",
)

// Service structure with the configuration for the API service.
type Service struct {
    Configuration Config
}

// Providers structure to facilitate the building process as a factory.
type Providers struct {
    networkProvider     networkstorage.Provider
    clusterProvider     clusterstorage.Provider
    nodeProvider        nodestorage.Provider
    appDescProvider     appdescstorage.Provider
    appInstProvider     appinststorage.Provider
    userProvider        userstorage.Provider
    accessProvider      accessstorage.Provider
    passwordProvider    passwordstorage.Provider
    oauthProvider       oauthstorage.Provider
    credentialsProvider credentialsstorage.Provider
    sessionProvider     sessionstorage.Provider
    configProvider      configstorage.Provider
}

// Name of the service.
func (s *Service) Name() string {
    return "System Model Service."
}

// Description of the service.
func (s *Service) Description() string {
    return "Api service of the System Model project."
}

// CreateInMemoryProviders returns a set of in-memory providers.
func (s *Service) CreateInMemoryProviders() *Providers {
    return &Providers{
        networkstorage.NewMockupNetworkProvider(),
        clusterstorage.NewMockupClusterProvider(),
        nodestorage.NewMockupNodeProvider(),
        appdescstorage.NewMockupAppDescProvider(),
        appinststorage.NewMockupAppInstProvider(),
        userstorage.NewMockupUserProvider(),
        accessstorage.NewMockupUserAccessProvider(),
        passwordstorage.NewMockupPasswordProvider(),
        oauthstorage.NewMockupOAuthProvider(),
        credentialsstorage.NewMockupCredentialsProvider(),
        sessionstorage.NewMockupSessionProvider(),
        configstorage.NewMockupConfigProvider()}
}

// CreateFileSystemProviders returns a set of filesystem-backed providers.
func (s *Service) CreateFileSystemProviders(basePath string) *Providers {
    return &Providers{
        networkstorage.NewFileSystemProvider(basePath),
        clusterstorage.NewFileSystemProvider(basePath),
        nodestorage.NewFileSystemProvider(basePath),
        appdescstorage.NewFileSystemProvider(basePath),
        appinststorage.NewFileSystemProvider(basePath),
        userstorage.NewFileSystemProvider(basePath),
        accessstorage.NewFileSystemProvider(basePath),
        passwordstorage.NewFileSystemProvider(basePath),
        oauthstorage.NewFileSystemProvider(basePath),
        credentialsstorage.NewFileSystemProvider(basePath),
        sessionstorage.NewFileSystemProvider(basePath),
        configstorage.NewFileSystemProvider(basePath)}
}

// GetProviders builds the providers according to the selected backend.
func (s *Service) GetProviders() *Providers {
    if s.Configuration.UseFileSystemProvider {
        return s.CreateFileSystemProviders(s.Configuration.FileSystemBasePath)
    } else if s.Configuration.UseInMemoryProviders {
        return s.CreateInMemoryProviders()
    } else {
        loggerService.Fatal("Unsupported type of provider")
    }
    return nil
}

// Run the service, launch the REST service handler.
func (s *Service) Run() error {
    loggerService.Info("Starting Api Service.")
    s.Configuration.Print()
    handler := mux.NewRouter()

    p := s.GetProviders()

    networkMgr := network.NewManager(p.networkProvider)
    clusterMgr := cluster.NewManager(p.networkProvider, p.clusterProvider, p.appInstProvider)
    nodeMgr := node.NewManager(p.networkProvider, p.clusterProvider, p.nodeProvider)

    networkHandler := network.NewHandler(networkMgr)
    networkHandler.SetRoutes(handler)

    clusterHandler := cluster.NewHandler(clusterMgr)
    clusterHandler.SetRoutes(handler)

    nodeHandler := node.NewHandler(nodeMgr)
    nodeHandler.SetRoutes(handler)

    appManager := app.NewManager(p.networkProvider, p.appDescProvider, p.appInstProvider)
    appHandler := app.NewHandler(appManager)
    appHandler.SetRoutes(handler)

    dumpManager := dump.NewManager(p.networkProvider, p.clusterProvider, p.nodeProvider,
        p.appDescProvider, p.appInstProvider, p.userProvider, p.accessProvider)
    dumpHandler := dump.NewHandler(dumpManager)
    dumpHandler.SetRoutes(handler)

    backupRestoreManager := backup.NewManager(p.networkProvider, p.clusterProvider, p.nodeProvider,
        p.appDescProvider, p.userProvider, p.accessProvider, p.passwordProvider)
    backupRestoreHandler := backup.NewHandler(backupRestoreManager)
    backupRestoreHandler.SetRoutes(handler)

    infoManager := info.NewManager(p.networkProvider, p.clusterProvider, p.nodeProvider,
        p.appDescProvider, p.appInstProvider, p.userProvider)
    infoHandler := info.NewHandler(infoManager)
    infoHandler.SetRoutes(handler)

    userManager := user.NewManager(p.userProvider, p.accessProvider, p.passwordProvider, p.oauthProvider)
    userHandler := user.NewHandler(userManager)
    userHandler.SetRoutes(handler)

    accessManager := access.NewManager(p.accessProvider)
    accessHandler := access.NewHandler(accessManager)
    accessHandler.SetRoutes(handler)

    passwordManager := password.NewManager(p.passwordProvider)
    passwordHandler := password.NewHandler(passwordManager)
    passwordHandler.SetRoutes(handler)

    oauthManager := oauth.NewManager(p.oauthProvider)
    oauthHandler := oauth.NewHandler(oauthManager)
    oauthHandler.SetRoutes(handler)

    sessionManager := session.NewManager(p.sessionProvider)
    sessionHandler := session.NewHandler(sessionManager)
    sessionHandler.SetRoutes(handler)

    credentialsManager := credentials.NewManager(p.credentialsProvider)
    credentialsHandler := credentials.NewHandler(credentialsManager)
    credentialsHandler.SetRoutes(handler)

    configManager := config.NewManager(p.configProvider)
    configHandler := config.NewHandler(configManager)
    configHandler.SetRoutes(handler)

    handler.HandleFunc("/ping", s.ping).Methods("GET")
    loggerService.Info("Ready to serve HTTP")
    s.listRoutes(handler)
    go http.ListenAndServe(fmt.Sprintf(":%d", s.Configuration.Port), handler)

    // Create default admin
    if !s.createDefaultAdmin(userManager, accessManager, passwordManager, oauthManager) {
        s.Finalize(true)
    }

    if !s.createDefaultConfig(configManager) {
        s.Finalize(true)
    }
    loggerService.Info("Do not forget to change the default user settings!!!!!!")

    return nil
}

func (s *Service) createDefaultConfig(configManager config.Manager) bool {
    config, _ := configManager.GetConfig()
    if config == nil {
        defaultConfig := entities.NewConfig("168h")
        loggerService.Infof("Creating default config ...")
        err := configManager.SetConfig(*defaultConfig)
        if err != nil {
            loggerService.Error(err)
            return false
        }
        loggerService.Info("Configuration stored")
    }
    return true
}

// Create a default admin if needed in the system model.
func (s *Service) createDefaultAdmin(userMg user.Manager, accessMg access.Manager, passwordMg password.Manager,
    oauthMg oauth.Manager) bool {
    // Create the new userMg entry
    returnedUser, err := userMg.GetUser(s.Configuration.DefaultAdminUser)
    if returnedUser == nil {
        if err.Error() == "[Operation] "+errors.UserDoesNotExist {
            loggerService.Infof("Creating default admin user %s...", s.Configuration.DefaultAdminUser)
            // We set the expiration day for the next year
            _, err := userMg.AddUser(*entities.NewAddUserRequest(s.Configuration.DefaultAdminUser,
                "admin", "admin", "admin@daisho.group", time.Now(), time.Now().Add(time.Hour*8760)))
            if err != nil {
                loggerService.Error(err)
                return false
            }
            loggerService.Info("Admin user created")
        } else {
            loggerService.Error("System model is not ready to create default admin")
            loggerService.Error(err.Error())
            loggerService.Error(errors.UserDoesNotExist)
            loggerService.Error(err.Error() == errors.UserDoesNotExist)

            return false
        }
    } else {
        loggerService.Infof("Default admin user %s is already defined", s.Configuration.DefaultAdminUser)
    }

    // Create the accessMg entry
    loggerService.Infof("Creating default admin user %s... privilege", s.Configuration.DefaultAdminUser)
    adminRole := []entities.RoleType{entities.GlobalAdmin}
    _, err = accessMg.AddAccess(s.Configuration.DefaultAdminUser, *entities.NewAddUserAccessRequest(adminRole))
    if err != nil {
        loggerService.Error(err)
        return false
    }
    loggerService.Info("Done")

    // Create the passwordMg
    loggerService.Infof("Creating default admin %s password...", s.Configuration.DefaultAdminUser)
    p, err := entities.NewPassword(s.Configuration.DefaultAdminUser, &s.Configuration.DefaultAdminPassword)
    if err != nil {
        loggerService.Error(err)
        return false
    }
    err = passwordMg.SetPassword(*p)
    if err != nil {
        loggerService.Error(err)
        return false
    }
    loggerService.Info("Done")

    // Create the OAuth entry
    // The secret must have already been set prior to run this command or the access request will fail.
    loggerService.Infof("Creating default admin %s oauth entry...", s.Configuration.DefaultAdminUser)

    err = oauthMg.SetSecret(s.Configuration.DefaultAdminUser, entities.NewOAuthAddEntryRequest("samurai",
        s.Configuration.DefaultAdminUser, s.Configuration.DefaultAdminPassword))
    if err != nil {
        loggerService.Error(err)
        return false
    }
    loggerService.Info("Done")

    return true
}

func (s *Service) listRoutes(handler *mux.Router) {
    handler.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
        path, err := route.GetPathTemplate()
        if err != nil {
            return err
        }
        methods, err := route.GetMethods()
        if err != nil {
            return err
        }
        methodStr := ""
        for _, m := range methods {
            methodStr = methodStr + ", " + m
        }
        log.Info(methodStr + " : " + path)
        return nil
    })
}

// The ping endpoint returns a 200 Ok response for external services to check that the service is up.
// TODO If used for liveness probes, enhance the check with provider status.
func (s *Service) ping(w http.ResponseWriter, r *http.Request) {
    dhttp.RespondWithJSON(w, http.StatusOK, entities.NewSuccessfulOperation("ping"))
}

// Finalize the service.
func (s *Service) Finalize(killSignal bool) {
    loggerService.Info("Finalize Api Service.")
}
