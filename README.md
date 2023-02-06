# Project Design

## Overall considerations

This project setup is based on the one advocated by Bill Kennedy of Ardan Labs in the Ultimate Go: Web Services 3.0 
course. I have attempted to build my WebMind framework based on this structure, and mention the essential aspects 
of the project layout here:

### Layers

The project is not broken up into groupings, but instead is all about layers: Each level in the hierarchy should 
have no more than five layers, in order to make the mental model simple enough for humans to comprehend. Imports should
NEVER be across folders on the same level, and only from top to bottom. This avoids circular imports and keeps the 
folders on one level of the hierarchy independent of one another. 

The top level is broken up into the following layers:

- App
- Business
- Foundation
- Vendor
- Zarf

The five layers mentioned here have their names chosen so they are alphabetically on top of one another.
These layers can have further sub-layerings, but try to stick to the five layer maximum, and if possible 
find names that again reflect their layeredness in the correct order.

### App layer

This layer is the topmost layer of the program, mainly concerned with starting and stopping the app, 
and performing input validations and output actions. These can be CLI, web, or other UI interfaces.

### Business layer

Imported by the app layer, business provides sub-layers for the more business-oriented layers of the project.
Usually a sub-layering into care, data, sys and web suffices here. An exception to the above import restrictions 
here is that Bill normally allows app layers to import the data layer directly rather than having to write many 
very thin access methods in core for them. 

#### Business/Core layer

Provides the API into the Business layer. It can be bypassed for trivial access to say the database layer, 
but mostly provides the more complex and structured ways of accessing business data. If an app requires multiple 
calls into the business/data layer to construct some output, the function doing that and providing it to 
the app layer would be in business/core.

#### Business/Data layer

This is usually subdivided into schema, store and possibly test. Also might have a metrics sub-layer to provide 
business-related metrics. Usually this is where the CRUD layers reside.

##### Business/Data/Store layer

Data store layer is where the business entities are. For instance in a webstore application you would typically find 
here the sub-layers product, customer, etc. 

### Business/Sys layer

Usually subdivided into auth, database, metrics and validate, this handles the various aspects of data provided to 
higher layers. But this is still at a level very much related to the 'business at hand'.

### Foundation layer

Harbouring functions that could potentially be re-used across multiple applications, this is where most of WebMind's
functionality resides for now. It could be sub-divided into levels like docker, keystore, logger, web and worker, 
but don't sub-divide until there is a need to. These packages form the 'standard library for the project'. They 
should not log, and can only return wrapped errors. Should be truly REUSEABLE!

### Vendor layer

This is where we store the vendored packages of third party packages we wish to use. Typical subfolders of this are 
named github.com, go.uber.org, etc. This makes sure we always have the code to third party packages available, no 
matter what happens to the original package.

### Zarf layer

Named after the sleeve around a hot cup of coffee, the Zarf layer will provide anything related to configuration, 
distribution, and deployment. This would typically contain stuff related to Docker, K8s, Kind, etc. 

#### Zarf/Docker layer

Rather than have one general dockerfile for all apps, Bill prefers to use one dockerfile per app, so it is easier 
to debug.  



