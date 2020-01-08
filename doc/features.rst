.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Features

.. _rs_server_features:

Key Features
============

Digital Rebar Provision is a new generation of data center automation
designed for operators with a cloud-first approach. Data center
provisioning is surprisingly complex because it’s caught between
cutting edge hardware and arcane protocols embedded in firmware
requirements that are still very much alive.

Swagger REST API & CLI
----------------------

Cloud-first means having a great, tested API. Years of provisioning
experience went into this 3rd generation design and it shows. That
includes a powerful API-driven DHCP server.

Security & Authenticated API
----------------------------

Not an afterthought, we use both HTTPS and user authentication for the
API. Our mix of basic and bearer token authentication recognizes that
both users and automation will use the API. This brings a new level of
security and control to data center provisioning.

Stand-alone multi-architecture Golang binary
--------------------------------------------

There are no dependencies or prerequisites, plus upgrades are drop in
replacements. That allows users to experiment isolated on their laptop
and then easily register it as a SystemD service.

Nested Template Expansion
-------------------------

In Digital Rebar Provision, Boot Environments are composed of reusable
template snippets. These templates can incorporate global, profile or
machine specific properties that enable users to set services, users,
security or scripting extensions for their environment.  Configuration
at Global, Group/Profile and Node level. Properties for templates can
be managed in a wide range of ways that allows operators to manage
large groups of servers in consistent ways.

See :ref:`rs_data_template` for about using template expansion variables.

Multi-mode (but optional) DHCP
------------------------------

Network IP allocation is a key component of any provisioning
infrastructure; however, DHCP needs are highly site dependent. Digital
Rebar Provision works as a multi-interface DHCP listener and can also
resolve addresses from DHCP forwarders. It can even be disabled if the
environment already has a DHCP service that can configure a the “next
boot” provider.

Dynamic Provisioner templates for TFTP and HTTP
-----------------------------------------------

For security and scale, Digital Rebar Provision builds provisioning
files dynamically based on the Boot Environment Template system. This
means that critical system information is not written to disk and
files do not have to be synchronized. Of course, when a file needs to
be served it works too.

Node Discovery Bootstrapping
----------------------------

Digital Rebar’s long-standing discovery process is enabled in the
Provisioner with the included discovery boot environment. That process
includes an integrated secure token sequence so that new machines can
self-register with the service via the API. This eliminates the need
to pre-populate the Digital Rebar Provision system.

Multiple Seeding Operating Systems
----------------------------------

Digital Rebar Provision comes with a long list of Boot Environments
and Templates including support for many Linux flavors, Windows, ESX
and even LinuxKit. Our template design makes it easy to expand and
update templates even on existing deployments.

Two-stage TFTP/HTTP Boot
------------------------

Our specialized Sledgehammer and Discovery images are designed for
speed with optimized install cycles the improve boot speed by
switching from PXE TFTP to IPXE HTTP in a two stage process. This
ensures maximum hardware compatibility without creating excess network
load.

Context Changing
----------------

Digital Rebar Workflows allow the operational location of the runner/agent to changed during a workflow.  By default, the runner works on the machine; however, v4.1 and later platforms are able to migrate the context into a container running on the endpoint.  This allows workflows to leverage multiple network and security trust zones.  It also allows them to perform work based on libraries managed in a container.