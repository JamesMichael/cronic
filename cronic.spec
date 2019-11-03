undefine _missing_build_ids_terminate_build

Name:          cronic
Summary:       Interactive cron invocation tool
Version:       1.1.0
Release:       1%{?dist}

Group:         Applications
License:       Public Domain
URL:           https://github.com/JamesMichael/cronic
Source0:       %{name}-%{version}-%{release}.tar.gz
BuildRoot:     %{_tmppath}/%{name}-%{version}-%{release}-root
BuildRequires: golang

%description
Interactive cron invocation tool

%prep

%setup -q -n cronic

%build

pushd cmd/cronic
go build -mod=vendor
./cronic --help-man | gzip -9 > cronic.1.gz

%install
install -D cmd/cronic/cronic      %{buildroot}/%{_bindir}/cronic
install -D cmd/cronic/cronic.1.gz %{buildroot}/%{_mandir}/man1/cronic.1.gz

%clean

rm -rf $RPM_BUILD_ROOT

%files

%defattr(-, root, root, -)
%{_bindir}/cronic
%{_mandir}/man1/cronic.1.gz

%changelog
* Sun Nov 03 2019 James Michael <jamesamichael@gmail.com> - 1.1.0-1
- Exit after running command
- Add -path flag, CRONIC_PATH environment variable

* Sun Jul 14 2019 James Michael <jamesamichael@gmail.com> - 1.0.1-1
- Ensure the selected command is run (#1)

* Sun May 19 2019 James Michael <jamesamichael@gmail.com> - 1.0.0-1
- Initial package
