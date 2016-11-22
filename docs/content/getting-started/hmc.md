---
date: 2016-11-21T15:04:24+01:00
title: HMC
menu:
  main:
    parent: Getting started
    identifier: /getting-started/hmc
    weight: 40
---

# Enabling PCM data collection

By default, PCM data are not collected. It's needed to enable it at the HMC level.

  {{< gallery image="hmc_pcm.png" >}}

You can enable PCM data collection by managed system on enable it for all of them.

{{< gallery image="hmc_pcm_settings.png" >}}

# firewall settings

To fetch PCM data from the HMC, you need to have the port **12443** allowed.
